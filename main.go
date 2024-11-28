package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

func LoadEnv() {
	err := godotenv.Load("migrations/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
func main() {
	//загружаем файл env
	LoadEnv()
	// Читаем переменные из окружения
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	applied, err := goose.Up(db, "./migrations")
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	log.Printf("%d Migrations applied successfully", applied)

	//запускаем сервер
	logrus.Info("Запускаем сервер")
	http.Handlefunc()
	err = http.ListenAndServe(":dbPort", nil)
	if err != nil {
		panic(err)
	}

}
