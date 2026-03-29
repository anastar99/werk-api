package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anastar99/werk-api/server"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

// ← THIS is the PostgreSQL driver

var db *sql.DB

func initDB() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables ")
	}

	// Get the database URL from env
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set in environment")
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB successfully")
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	initDB()
	s := &server.Server{
		DB: db,
	}

	log.Println("Starting out simple server...")

	s.Routes()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080" // defualt if PORT not set
	}

	log.Println("Started on port", port)
	fmt.Println("To close connection CTRL+C")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
