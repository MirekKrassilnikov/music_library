package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/MirekKrassilnikov/music_library/handlers"
	"github.com/MirekKrassilnikov/music_library/services"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

func LoadEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
func main() {
	//загружаем файл env
	LoadEnv()
	// Читаем переменные из окружения
	//dbHost := os.Getenv("DB_HOST")
	//dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	//dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connectStr := fmt.Sprintf("user=%s dbname=postgres sslmode=disable", dbUser)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)).Scan(&exists)
	if err != nil {
		log.Fatalf("failed to check if database exists: %v", err)
	}

	if !exists {
		// Создаем базу данных, если она не существует
		_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.Fatalf("failed to create database: %v", err)
		}
		log.Printf("Database %s created successfully", dbName)
	} else {
		log.Printf("Database %s already exists", dbName)
	}

	err = goose.Up(db, "./migrations")
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	log.Printf("%d Migrations applied successfully")
	// Создаем SongService
	songService := &services.SongService{DB: db}

	// Создаем SongHandler
	songHandler := &handlers.SongHandler{SongService: songService}

	//запускаем сервер
	logrus.Info("Запускаем сервер")
	http.HandleFunc("/songs", songHandler.HandleGetAllSongs)
	http.HandleFunc("/lyrics", songHandler.HandleGetLyrics)
	http.HandleFunc("/delete", songHandler.HandleDeleteSong)
	http.HandleFunc("/add", songHandler.HandleAddNewSong)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
