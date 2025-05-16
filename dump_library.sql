--
-- PostgreSQL database dump
--

-- Dumped from database version 17.3 (Debian 17.3-1.pgdg120+1)
-- Dumped by pg_dump version 17.3 (Debian 17.3-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: library; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA library;


ALTER SCHEMA library OWNER TO postgres;

--
-- Name: decrease_available_copies(); Type: FUNCTION; Schema: library; Owner: postgres
--

CREATE FUNCTION library.decrease_available_copies() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  UPDATE "library".books
  SET available_copies = available_copies - 1
  WHERE book_id = NEW.book_id;

  RETURN NEW;
END;
$$;


ALTER FUNCTION library.decrease_available_copies() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: authors; Type: TABLE; Schema: library; Owner: postgres
--

CREATE TABLE library.authors (
    author_id integer NOT NULL,
    name character varying(250) NOT NULL,
    date_of_birth date,
    country character varying(200),
    bio text,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT check_author_name CHECK (((name)::text ~ '^[А-ЯЁ]\.\s[А-ЯЁ]\.\s[А-ЯЁа-яё-]+$'::text)),
    CONSTRAINT check_date_of_birth CHECK ((date_of_birth < CURRENT_DATE))
);


ALTER TABLE library.authors OWNER TO postgres;

--
-- Name: authors_author_id_seq; Type: SEQUENCE; Schema: library; Owner: postgres
--

CREATE SEQUENCE library.authors_author_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE library.authors_author_id_seq OWNER TO postgres;

--
-- Name: authors_author_id_seq; Type: SEQUENCE OWNED BY; Schema: library; Owner: postgres
--

ALTER SEQUENCE library.authors_author_id_seq OWNED BY library.authors.author_id;


--
-- Name: book_authors; Type: TABLE; Schema: library; Owner: postgres
--

CREATE TABLE library.book_authors (
    book_id integer NOT NULL,
    author_id integer NOT NULL
);


ALTER TABLE library.book_authors OWNER TO postgres;

--
-- Name: books; Type: TABLE; Schema: library; Owner: postgres
--

CREATE TABLE library.books (
    book_id integer NOT NULL,
    title character varying(250) NOT NULL,
    genre character varying(100),
    isbn character varying(20),
    total_copies integer NOT NULL,
    available_copies integer NOT NULL,
    published_date date,
    publisher character varying(250),
    description text,
    cover_image text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT books_available_copies_check CHECK ((available_copies >= 0)),
    CONSTRAINT books_total_copies_check CHECK ((total_copies > 0)),
    CONSTRAINT check_published_date CHECK ((published_date < CURRENT_DATE))
);


ALTER TABLE library.books OWNER TO postgres;

--
-- Name: books_book_id_seq; Type: SEQUENCE; Schema: library; Owner: postgres
--

CREATE SEQUENCE library.books_book_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE library.books_book_id_seq OWNER TO postgres;

--
-- Name: books_book_id_seq; Type: SEQUENCE OWNED BY; Schema: library; Owner: postgres
--

ALTER SEQUENCE library.books_book_id_seq OWNED BY library.books.book_id;


--
-- Name: borrowings; Type: TABLE; Schema: library; Owner: postgres
--

CREATE TABLE library.borrowings (
    borrowing_id integer NOT NULL,
    user_id integer,
    book_id integer,
    borrow_date date DEFAULT now(),
    return_date date,
    status character varying(50) DEFAULT 'borrowed'::character varying NOT NULL,
    due_date date GENERATED ALWAYS AS ((borrow_date + 14)) STORED,
    CONSTRAINT borrowings_status_check CHECK (((status)::text = ANY ((ARRAY['borrowed'::character varying, 'returned'::character varying, 'overdue'::character varying])::text[])))
);


ALTER TABLE library.borrowings OWNER TO postgres;

--
-- Name: borrowings_borrowing_id_seq; Type: SEQUENCE; Schema: library; Owner: postgres
--

CREATE SEQUENCE library.borrowings_borrowing_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE library.borrowings_borrowing_id_seq OWNER TO postgres;

--
-- Name: borrowings_borrowing_id_seq; Type: SEQUENCE OWNED BY; Schema: library; Owner: postgres
--

ALTER SEQUENCE library.borrowings_borrowing_id_seq OWNED BY library.borrowings.borrowing_id;


--
-- Name: penalties; Type: TABLE; Schema: library; Owner: postgres
--

CREATE TABLE library.penalties (
    penalty_id integer NOT NULL,
    user_id integer,
    amount numeric(10,2) NOT NULL,
    reason text,
    paid boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT penalties_amount_check CHECK ((amount >= (0)::numeric))
);


ALTER TABLE library.penalties OWNER TO postgres;

--
-- Name: penalties_penalty_id_seq; Type: SEQUENCE; Schema: library; Owner: postgres
--

CREATE SEQUENCE library.penalties_penalty_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE library.penalties_penalty_id_seq OWNER TO postgres;

--
-- Name: penalties_penalty_id_seq; Type: SEQUENCE OWNED BY; Schema: library; Owner: postgres
--

ALTER SEQUENCE library.penalties_penalty_id_seq OWNED BY library.penalties.penalty_id;


--
-- Name: reservations; Type: TABLE; Schema: library; Owner: postgres
--

CREATE TABLE library.reservations (
    reserve_id integer NOT NULL,
    user_id integer,
    book_id integer,
    reservation_date date DEFAULT now(),
    status character varying(50) DEFAULT 'active'::character varying NOT NULL,
    end_of_reserve date GENERATED ALWAYS AS ((reservation_date + 3)) STORED,
    CONSTRAINT reservations_status_check CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'canceled'::character varying, 'expired'::character varying])::text[])))
);


ALTER TABLE library.reservations OWNER TO postgres;

--
-- Name: reservations_reserve_id_seq; Type: SEQUENCE; Schema: library; Owner: postgres
--

CREATE SEQUENCE library.reservations_reserve_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE library.reservations_reserve_id_seq OWNER TO postgres;

--
-- Name: reservations_reserve_id_seq; Type: SEQUENCE OWNED BY; Schema: library; Owner: postgres
--

ALTER SEQUENCE library.reservations_reserve_id_seq OWNED BY library.reservations.reserve_id;


--
-- Name: users; Type: TABLE; Schema: library; Owner: postgres
--

CREATE TABLE library.users (
    user_id integer NOT NULL,
    login character varying(100) NOT NULL,
    email character varying(100) NOT NULL,
    password_hash text NOT NULL,
    first_name text NOT NULL,
    surname text NOT NULL,
    patronymic text,
    date_of_birth date,
    phone character varying(20),
    role text DEFAULT 'user'::text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT check_date_of_birth CHECK ((date_of_birth < CURRENT_DATE)),
    CONSTRAINT check_first_name CHECK ((first_name ~ '^[A-ZА-ЯЁ][a-zа-яё]+$'::text)),
    CONSTRAINT check_patronymic CHECK ((patronymic ~ '^[A-ZА-ЯЁ][a-zа-яё]+$'::text)),
    CONSTRAINT check_surname CHECK ((surname ~ '^[A-ZА-ЯЁ][a-zа-яё]+$'::text)),
    CONSTRAINT users_role_check CHECK ((role = ANY (ARRAY['user'::text, 'admin'::text])))
);


ALTER TABLE library.users OWNER TO postgres;

--
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: library; Owner: postgres
--

CREATE SEQUENCE library.users_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE library.users_user_id_seq OWNER TO postgres;

--
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: library; Owner: postgres
--

ALTER SEQUENCE library.users_user_id_seq OWNED BY library.users.user_id;


--
-- Name: messages; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.messages (
    id bigint NOT NULL,
    text text
);


ALTER TABLE public.messages OWNER TO postgres;

--
-- Name: messages_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.messages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.messages_id_seq OWNER TO postgres;

--
-- Name: messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.messages_id_seq OWNED BY public.messages.id;


--
-- Name: authors author_id; Type: DEFAULT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.authors ALTER COLUMN author_id SET DEFAULT nextval('library.authors_author_id_seq'::regclass);


--
-- Name: books book_id; Type: DEFAULT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.books ALTER COLUMN book_id SET DEFAULT nextval('library.books_book_id_seq'::regclass);


--
-- Name: borrowings borrowing_id; Type: DEFAULT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.borrowings ALTER COLUMN borrowing_id SET DEFAULT nextval('library.borrowings_borrowing_id_seq'::regclass);


--
-- Name: penalties penalty_id; Type: DEFAULT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.penalties ALTER COLUMN penalty_id SET DEFAULT nextval('library.penalties_penalty_id_seq'::regclass);


--
-- Name: reservations reserve_id; Type: DEFAULT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.reservations ALTER COLUMN reserve_id SET DEFAULT nextval('library.reservations_reserve_id_seq'::regclass);


--
-- Name: users user_id; Type: DEFAULT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.users ALTER COLUMN user_id SET DEFAULT nextval('library.users_user_id_seq'::regclass);


--
-- Name: messages id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.messages ALTER COLUMN id SET DEFAULT nextval('public.messages_id_seq'::regclass);


--
-- Data for Name: authors; Type: TABLE DATA; Schema: library; Owner: postgres
--

COPY library.authors (author_id, name, date_of_birth, country, bio, created_at) FROM stdin;
7	А. С. Пушкин	1799-05-26	Россия	Русский писатель, один из величайших авторов мировой литературы.	2025-04-03 08:54:15.998244
8	Л. Н. Толстой	1828-09-09	Россия	Русский писатель, мыслитель и публицист, автор эпопеи 'Война и мир' и романа 'Анна Каренина'. Один из величайших классиков мировой литературы.	2025-04-10 16:14:29.739038
\.


--
-- Data for Name: book_authors; Type: TABLE DATA; Schema: library; Owner: postgres
--

COPY library.book_authors (book_id, author_id) FROM stdin;
31	7
33	8
\.


--
-- Data for Name: books; Type: TABLE DATA; Schema: library; Owner: postgres
--

COPY library.books (book_id, title, genre, isbn, total_copies, available_copies, published_date, publisher, description, cover_image, created_at, updated_at) FROM stdin;
33	Война и мир	Исторический роман	978-5170708573	50	50	1869-01-01	Издательство Эксмо	«Война и мир» — эпический роман Льва Толстого, охватывающий судьбы нескольких дворянских семей на фоне Отечественной войны 1812 года. Произведение сочетает философские размышления, исторические события и личные переживания героев.	https://example.com/war-and-peace-cover.jpg	2025-04-10 16:16:54.842473	2025-04-10 16:16:54.842473
31	Евгений Онегин	Роман в стихах	978-5171202414	30	29	1833-01-01	Издательство Эксмо	«Евгений Онегин» — это знаменитый роман в стихах Александра Пушкина, повествующий о судьбе светского молодого человека Онегина, его дружбе с Ленским и несостоявшейся любви с Татьяной Лариной.	https://example.com/evgeny-onegin-cover.jpg	2025-04-03 08:55:53.981716	2025-04-03 08:55:53.981716
\.


--
-- Data for Name: borrowings; Type: TABLE DATA; Schema: library; Owner: postgres
--

COPY library.borrowings (borrowing_id, user_id, book_id, borrow_date, return_date, status) FROM stdin;
\.


--
-- Data for Name: penalties; Type: TABLE DATA; Schema: library; Owner: postgres
--

COPY library.penalties (penalty_id, user_id, amount, reason, paid, created_at) FROM stdin;
\.


--
-- Data for Name: reservations; Type: TABLE DATA; Schema: library; Owner: postgres
--

COPY library.reservations (reserve_id, user_id, book_id, reservation_date, status) FROM stdin;
1	2	31	2025-04-03	active
3	2	33	2025-04-11	active
4	3	31	2025-04-11	active
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: library; Owner: postgres
--

COPY library.users (user_id, login, email, password_hash, first_name, surname, patronymic, date_of_birth, phone, role, created_at, updated_at) FROM stdin;
2	johndoe	johndoe@example.com	$2a$10$Ish04rnFpAk2hYUb4cdEmui8QqKxPB9ZtVjsNUUgP2rVhnQPbBohK	John	Doe	Michael	1990-05-15	+1234567890	user	2025-04-03 08:51:58.106104	2025-04-03 08:51:58.106104
3	annasmith	anna.smith@example.com	$2a$10$yrfBbtIwoGjZKqEcwxfneOZd6ulh0O7j3PcSGAJ33/qqgnxLOkpZm	Anna	Smith	Elena	1987-11-23	+79876543210	admin	2025-04-10 16:12:39.800188	2025-04-10 16:12:39.800188
\.


--
-- Data for Name: messages; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.messages (id, text) FROM stdin;
2	test
3	test
4	test
5	test
\.


--
-- Name: authors_author_id_seq; Type: SEQUENCE SET; Schema: library; Owner: postgres
--

SELECT pg_catalog.setval('library.authors_author_id_seq', 8, true);


--
-- Name: books_book_id_seq; Type: SEQUENCE SET; Schema: library; Owner: postgres
--

SELECT pg_catalog.setval('library.books_book_id_seq', 33, true);


--
-- Name: borrowings_borrowing_id_seq; Type: SEQUENCE SET; Schema: library; Owner: postgres
--

SELECT pg_catalog.setval('library.borrowings_borrowing_id_seq', 3, true);


--
-- Name: penalties_penalty_id_seq; Type: SEQUENCE SET; Schema: library; Owner: postgres
--

SELECT pg_catalog.setval('library.penalties_penalty_id_seq', 1, false);


--
-- Name: reservations_reserve_id_seq; Type: SEQUENCE SET; Schema: library; Owner: postgres
--

SELECT pg_catalog.setval('library.reservations_reserve_id_seq', 4, true);


--
-- Name: users_user_id_seq; Type: SEQUENCE SET; Schema: library; Owner: postgres
--

SELECT pg_catalog.setval('library.users_user_id_seq', 3, true);


--
-- Name: messages_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.messages_id_seq', 5, true);


--
-- Name: authors authors_name_key; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.authors
    ADD CONSTRAINT authors_name_key UNIQUE (name);


--
-- Name: authors authors_pkey; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.authors
    ADD CONSTRAINT authors_pkey PRIMARY KEY (author_id);


--
-- Name: book_authors book_authors_pkey; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.book_authors
    ADD CONSTRAINT book_authors_pkey PRIMARY KEY (book_id, author_id);


--
-- Name: books books_isbn_key; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.books
    ADD CONSTRAINT books_isbn_key UNIQUE (isbn);


--
-- Name: books books_pkey; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (book_id);


--
-- Name: borrowings borrowings_pkey; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.borrowings
    ADD CONSTRAINT borrowings_pkey PRIMARY KEY (borrowing_id);


--
-- Name: penalties penalties_pkey; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.penalties
    ADD CONSTRAINT penalties_pkey PRIMARY KEY (penalty_id);


--
-- Name: reservations reservations_pkey; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.reservations
    ADD CONSTRAINT reservations_pkey PRIMARY KEY (reserve_id);


--
-- Name: authors unique_name; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.authors
    ADD CONSTRAINT unique_name UNIQUE (name);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_login_key; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.users
    ADD CONSTRAINT users_login_key UNIQUE (login);


--
-- Name: users users_phone_key; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.users
    ADD CONSTRAINT users_phone_key UNIQUE (phone);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);


--
-- Name: reservations trg_decrease_available_copies; Type: TRIGGER; Schema: library; Owner: postgres
--

CREATE TRIGGER trg_decrease_available_copies AFTER INSERT ON library.reservations FOR EACH ROW EXECUTE FUNCTION library.decrease_available_copies();


--
-- Name: book_authors book_authors_author_id_fkey; Type: FK CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.book_authors
    ADD CONSTRAINT book_authors_author_id_fkey FOREIGN KEY (author_id) REFERENCES library.authors(author_id) ON DELETE CASCADE;


--
-- Name: book_authors book_authors_book_id_fkey; Type: FK CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.book_authors
    ADD CONSTRAINT book_authors_book_id_fkey FOREIGN KEY (book_id) REFERENCES library.books(book_id) ON DELETE CASCADE;


--
-- Name: borrowings borrowings_book_id_fkey; Type: FK CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.borrowings
    ADD CONSTRAINT borrowings_book_id_fkey FOREIGN KEY (book_id) REFERENCES library.books(book_id) ON DELETE CASCADE;


--
-- Name: borrowings borrowings_user_id_fkey; Type: FK CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.borrowings
    ADD CONSTRAINT borrowings_user_id_fkey FOREIGN KEY (user_id) REFERENCES library.users(user_id) ON DELETE CASCADE;


--
-- Name: penalties penalties_user_id_fkey; Type: FK CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.penalties
    ADD CONSTRAINT penalties_user_id_fkey FOREIGN KEY (user_id) REFERENCES library.users(user_id) ON DELETE CASCADE;


--
-- Name: reservations reservations_book_id_fkey; Type: FK CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.reservations
    ADD CONSTRAINT reservations_book_id_fkey FOREIGN KEY (book_id) REFERENCES library.books(book_id) ON DELETE CASCADE;


--
-- Name: reservations reservations_user_id_fkey; Type: FK CONSTRAINT; Schema: library; Owner: postgres
--

ALTER TABLE ONLY library.reservations
    ADD CONSTRAINT reservations_user_id_fkey FOREIGN KEY (user_id) REFERENCES library.users(user_id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

