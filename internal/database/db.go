package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var DB *pgx.Conn

func ConnectDB() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	url := fmt.Sprintf("postgresql://%s:%s@localhost:5432/tinder?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))

	DB, err = pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	fmt.Println("✅ Подключено к PostgreSQL")
}
