package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/MirekKrassilnikov/music_library/routes"
	_ "github.com/lib/pq"

	"github.com/MirekKrassilnikov/music_library/domain/services"
	"github.com/MirekKrassilnikov/music_library/handlers"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

//func to load .env file
func LoadEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
func main() {
	handler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(handler)
	//loading env file
	LoadEnv()
	// reading variables
	dbUser := os.Getenv("DB_USER")
	//dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	//connecting to default database
	connectStr := fmt.Sprintf("user=%s dbname=postgres sslmode=disable", dbUser)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	//Check if database dbName exists
	var exists bool
	err = db.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)).Scan(&exists)
	if err != nil {
		log.Fatalf("failed to check if database exists: %v", err)
	}

	if !exists {
		// Createing database dbName
		_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.Fatalf("failed to create database: %v", err)
		}
		logger.Info("Database %s created successfully", dbName)
	} else {
		log.Printf("Database %s already exists", dbName)
	}

	//Migration
	err = goose.Up(db, "./migrations")
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	logger.Info("%d Migrations applied successfully")

	// Create SongService
	songService := &services.SongService{DB: db}

	// Create SongHandler
	songHandler := &handlers.SongHandler{SongService: songService}

	// Create routes
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, songHandler)

	// Starting server
	logger.Info("Starting server on :8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.Error("Server failed to start", slog.String("error", err.Error()))
	}
	
}
