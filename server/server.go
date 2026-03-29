package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	DB *sql.DB
}

func (s *Server) ClockIn(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	query := `
		INSERT INTO
			attendance (day, clock_in) VALUES ($1, $2)
		`
	_, err := s.DB.Exec(query, today, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) ClockOut(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	query := `UPDATEattendance SET clock_out = $1 WHERE day = $2`
	_, err := s.DB.Exec(query, today, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (s *Server) WeeklyHours(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
	}

	rows, err := s.DB.Query(`
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

		if err := rows.Scan(&e.ID, &e.Day, &e.ClockIn, &e.ClockOut); err != nil {
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

func (s *Server) BiWeeklyHours(w http.ResponseWriter, r *http.Request) {
	// this reutnr the bikely summar of hours worked

	fmt.Println("Bi weekly hours")
}

func (s *Server) Entries(w http.ResponseWriter, r *http.Request) {

	rows, err := s.DB.Query(`
		SELECT
			id, day, clock_in, clock_out
		FROM
			attendance
	`)

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
		if err := rows.Scan(&e.ID, &e.Day, &e.ClockIn, &e.ClockOut); err != nil {
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
