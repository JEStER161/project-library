package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// Подключение к базе данных
func ConnectDB() {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=1234 sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	DB = pool
	log.Println("Подключение к БД установлено!")
}
