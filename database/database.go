package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	Id int64
	Password  string
	Site      string
	Notes     string
	Timestamp time.Time
	Nonce string
}

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "database/sanctum.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS entries (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			password TEXT NOT NULL,
			site TEXT NOT NULL,
			notes TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			nonce TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS master_password (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			password_hash TEXT NOT NULL,
			salt TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS auth_token (
			token_hash TEXT NOT NULL
		);`,
	}

	for i, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			log.Fatalf("Failed to execute query #%d: %v\nQuery: %s", i+1, err, query)
		}
	}
}

func GetMasterPassword() (string, string, error) {
	query := `SELECT password_hash, salt from master_password`

	row := DB.QueryRow(query)

	var	password_hash string
	var salt string
	err := row.Scan(&password_hash, &salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", nil
		}
		log.Fatalf("Error getting master password: %q: %s\n", err, query) 
	}
	return password_hash, salt,  nil
}

func InserMasterPassword(password string, salt string) (int64, error) {
	query := `
	INSERT INTO master_password (password_hash, salt)
	VALUES (?, ?);`

	result, err := DB.Exec(query, password, salt)
	if err != nil {
		log.Fatalf("Error inserting master password: %q: %s\n", err, query) 
	}
	return result.RowsAffected()
}

func GetEntries() ([]Entry, error) {
	query := `SELECT id, password, site, notes, timestamp FROM entries`

	rows, err := DB.Query(query)
	if err != nil {
		log.Fatalf("Error getting passwords: %q: %s\n", err, query) 
	}

	defer rows.Close()

	entries := []Entry{}
    for rows.Next() {
            var entry Entry
            err := rows.Scan(&entry.Id, &entry.Password, &entry.Site, &entry.Notes, &entry.Timestamp)
            if err != nil {
				log.Fatalf("Error getting passwords: %q: %s\n", err, query) 
            }
            entries = append(entries, entry)
    }
	return entries, nil
}


func InsertEntry(entry Entry) (int64, error) {
	query := `
	INSERT INTO entries (password, site, notes, timestamp, nonce)
	VALUES (?, ?, ?, ?, ?);`

	result, err := DB.Exec(query, entry.Password, entry.Site, entry.Notes, entry.Timestamp, entry.Nonce)
	if err != nil {
		log.Fatalf("Error inserting entry: %q: %s\n", err, query) 
	}
	return result.LastInsertId()
}

func DeleteEntry(id string) (error) {
	query := `DELETE FROM entries WHERE id = (?);`

	_, err := DB.Exec(query, id)
	if err != nil {
		log.Fatalf("Error deleting entry: %q: %s\n", err, query) 
	}
	return err 
}

func GetEntry(id string) (Entry, error) {
	query := `SELECT * FROM entries WHERE id = (?);`

	row := DB.QueryRow(query, id)

	entry := Entry{}
	err := row.Scan(&entry.Id, &entry.Password, &entry.Site, &entry.Notes, &entry.Timestamp, &entry.Nonce)
	if err != nil {
		return Entry{}, err
	}

	return entry, nil
}


func GetToken() (string, error) {
	query := `SELECT token_hash from auth_token`

	row := DB.QueryRow(query)

	var	token_hash string
	err := row.Scan(&token_hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		log.Fatalf("Error getting token: %q: %s\n", err, query) 
	}
	return token_hash, nil
}

func InsertToken(hashed_token string) (int64, error) {
	query := `
	INSERT INTO auth_token (token_hash)
	VALUES (?);`

	result, err := DB.Exec(query, hashed_token)
	if err != nil {
		log.Fatalf("Error inserting token: %q: %s\n", err, query) 
	}
	return result.RowsAffected()
}

func DeleteToken() (int64, error) {
	query := `DELETE from auth_token;`

	result, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Error inserting token: %q: %s\n", err, query) 
	}
	return result.RowsAffected()
}