package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// --------------------------------------------------------
//					STORE AND DATABASE
//---------------------------------------------------------

// store and database code exists in this file and can function
// independantly of other parts of the project. This separation
// lets this code get easily reused by both the cli tooling and
// the http web hosting code, since the database and store code
// does not need to know anything at all about the caller

type Store struct {
	db *sql.DB
}

type Entry struct {
	ID        int
	VisitedAt time.Time
	Place     string
	Comment   string
}

func init() {
	_ = godotenv.Load()
}

func newStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func mustOpenStore() *Store {

	db, err := sql.Open("mysql", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB connection failed:", err)
	}

	fmt.Println("Connected to MariaDB successfully")

	return newStore(db)
}

func (s *Store) ListEntries() ([]Entry, error) {

	rows, err := s.db.Query("SELECT id, visited_at, place, comment FROM entries ORDER BY visited_at ASC")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry

	for rows.Next() {
		var e Entry
		err := rows.Scan(&e.ID, &e.VisitedAt, &e.Place, &e.Comment)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}

func (s *Store) AddEntry(time time.Time, place string, comment string) error {

	_, err := s.db.Exec(`
        INSERT INTO entries 
		(visited_at, place, comment)
        VALUES (?, ?, ?)
    `, time, place, comment)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) RecentEntries(count int) ([]Entry, error) {

	rows, err := s.db.Query("SELECT id, visited_at, place, comment FROM entries ORDER BY visited_at DESC LIMIT ?", count)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry

	for rows.Next() {
		var e Entry
		err := rows.Scan(&e.ID, &e.VisitedAt, &e.Place, &e.Comment)
		if err != nil {
			return nil, err
		}

		entries = append(entries, e)
	}

	return entries, nil
}

func (s *Store) DeleteEntry(id int) error {
	_, err := s.db.Exec(`
		DELETE FROM entries
		WHERE id = ?
	`, id)

	return err
}

func (s *Store) UpdateEntry(id int, time time.Time, place string, comment string) error {
	_, err := s.db.Exec(`UPDATE entries 
	SET visited_at = ?, place = ?, comment = ? 
	WHERE id = ?
	`, time, place, comment, id)

	return err
}

func (s *Store) GetEntry(id int) (Entry, error) {
	var e Entry

	err := s.db.QueryRow(`SELECT id, visited_at, place, comment
	FROM entries
	WHERE id = ?`, id).Scan(&e.ID, &e.VisitedAt, &e.Place, &e.Comment)

	return e, err
}
