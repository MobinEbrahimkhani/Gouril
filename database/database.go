package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "gouril.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func Init(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			short_code TEXT UNIQUE NOT NULL,
			original_url TEXT NOT NULL
		);
	`

	_, err := db.Exec(query)
	return err
}

func InsertURL(db *sql.DB, shortCode string, originalURL string) error {
	query := `
		INSERT INTO urls (short_code, original_url)
		VALUES (?, ?);
	`

	_, err := db.Exec(query, shortCode, originalURL)
	return err
}

func GetURL(db *sql.DB, shortCode string) (string, error) {
	var originalURL string

	query := `SELECT original_url FROM urls WHERE short_code = ?`

	err := db.QueryRow(query, shortCode).Scan(&originalURL)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func ShortCodeExists(db *sql.DB, shortCode string) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS(
			SELECT 1 FROM urls WHERE short_code = ?
		)
	`

	err := db.QueryRow(query, shortCode).Scan(&exists)
	return exists, err
}
