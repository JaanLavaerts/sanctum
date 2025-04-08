package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/JaanLavaerts/sanctum/crypto"
	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	Password  string
	Site      string
	Notes     string
	Timestamp time.Time
}

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "database/sanctum.db")

	if err != nil {
		log.Fatal(err)
	}

	entriesQuery := `
 		CREATE TABLE IF NOT EXISTS entries (
  		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  		password TEXT NOT NULL,
		site TEXT NOT NULL,
		notes TEXT NOT NULL,
		timestamp DATETIME NOT NULL
	);`

	masterPasswordQuery := `
 		CREATE TABLE IF NOT EXISTS master_password (
  		password_hash TEXT NOT NULL
	);`


	 _, err = DB.Exec(entriesQuery)
 		if err != nil {
			log.Fatalf("Error creating table: %q: %s\n", err, entriesQuery) 
 		}

	 _, err = DB.Exec(masterPasswordQuery)
 		if err != nil {
			log.Fatalf("Error creating table: %q: %s\n", err, entriesQuery) 
 		}
}

func GetMasterPassword() (string, error) {
	query := `SELECT password_hash from master_password`

	row := DB.QueryRow(query)

	var	password_hash string
	err := row.Scan(&password_hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		log.Fatalf("Error getting master password: %q: %s\n", err, query) 
	}
	return password_hash, nil
}

func InserMasterPassword(plain_password string) (int64, error) {
	query := `
	INSERT INTO master_password (password_hash)
	VALUES (?);`

	hashed_password, err := crypto.GenerateHash(plain_password)
	if err != nil {
		log.Fatalf("Error creating hash: %q", err) 
	}
	result, err := DB.Exec(query, hashed_password)
	if err != nil {
		log.Fatalf("Error inserting master password: %q: %s\n", err, query) 
	}
	return result.RowsAffected()
}

func GetEntries() ([]Entry, error) {
	query := `SELECT password, site, notes, timestamp from entries`

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

