package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // ← THIS is the PostgreSQL driver
)

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

// Handler functions
func ClockIn(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	today := now.Format("2006-01-02") // YYYY-MM-DD

	query := `INSERT INTO attendance (day, clock_in) VALUES ($1, $2)`
	_, err := db.Exec(query, today, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated) // 201 Created
}

func ClockOut(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	today := now.Format("2006-01-02") // YYYY-MM-DD

	query := `UPDATE attendance SET clock_out = $1 WHERE day = $2`
	_, err := db.Exec(query, today, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func WeeklyHours(w http.ResponseWriter, r *http.Request) {

	// this returns the amount of hours that have been worked this week
	// in a list

	// Day 1: 3hrs
	// Day 2: 3.5hrs
	// Day 3: 4hrs
	// etc

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`
	SELECT 
		id, day, clock_in, clock_out
	FROM 
		attendance
	WHERE 
		day >= date_trunc('week', CURRENT_DATE AT TIME ZONE 'America/Los_Angeles')::date
	AND 
		day <= (CURRENT_DATE AT TIME ZONE 'America/Los_Angeles')::date
	ORDER BY
		day ASC;
`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	type Entry struct {
		ID       int        `json:"id"`
		Day      string     `json:"day"`
		Clockin  time.Time  `json:"clock_in"`
		ClockOut *time.Time `json:"clock_out,omitempty"`
	}

	var entries []Entry

	for rows.Next() {
		var e Entry
		var clockOut sql.NullTime

		if err := rows.Scan(&e.ID, &e.Day, &e.Clockin, &clockOut); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if clockOut.Valid {
			e.ClockOut = &clockOut.Time
		}

		entries = append(entries, e)
	}

	w.Header().Set("Content-Tyoe", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func BiWeeklyHours(w http.ResponseWriter, r *http.Request) {

	// this return the biweekly summary of hours worked

	// Day 1: 3.5hrs
	// Day 2: 4hrs
	// Day 3: 4hrs
	// ...
	// Total hours worked week 1

	// Day 1: 3.5hrs
	// ...
	// Total hours worked week 2

	fmt.Println("Bi weekly hours")
}

func Entries(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query(`SELECT id, day, clock_in, clock_out FROM attendance`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	type Entry struct {
		ID       int        `json:"id"`
		Day      string     `json:"day"`
		ClockIn  time.Time  `json:"clock_in"`
		ClockOut *time.Time `json:"clock_out,omitempty"`
	}

	var entries []Entry

	for rows.Next() {
		var e Entry
		var clockOut sql.NullTime
		if err := rows.Scan(&e.ID, &e.Day, &e.ClockIn, &clockOut); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if clockOut.Valid {
			e.ClockOut = &clockOut.Time
		}

		entries = append(entries, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)

}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	initDB()
	log.Println("Starting out simple server...")

	http.HandleFunc("/clock-in", ClockIn)
	http.HandleFunc("/clock-out", ClockOut)
	http.HandleFunc("/entries", Entries)
	http.HandleFunc("/week", WeeklyHours)
	http.HandleFunc("/biweek", BiWeeklyHours)

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
