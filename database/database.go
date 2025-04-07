package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	Password string
	Site string
	Notes string
	Timestamp time.Time
}

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "database/sanctum.db")

	if err != nil {
		log.Fatal(err)
	}

	query := `
 		CREATE TABLE IF NOT EXISTS entries (
  		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  		password TEXT NOT NULL,
		site TEXT NOT NULL,
		notes TEXT NOT NULL,
		timestamp DATETIME NOT NULL
	);`

	 _, err = DB.Exec(query)
 		if err != nil {
			log.Fatalf("Error creating table: %q: %s\n", err, query) 
 		}
}

func GetEntries() ([]Entry, error) {
	query := `SELECT password, site, notes, timestamp from entries LIMIT 1`

	rows, err := DB.Query(query)
	if err != nil {
		log.Fatalf("Error getting passwords: %q: %s\n", err, query) 
	}

	defer rows.Close()

	entries := []Entry{}
    for rows.Next() {
            var entry Entry
            err := rows.Scan(&entry.Password, &entry.Site, &entry.Notes, &entry.Timestamp)
            if err != nil {
				log.Fatalf("Error getting passwords: %q: %s\n", err, query) 
            }
            entries = append(entries, entry)
    }
	return entries, nil
}