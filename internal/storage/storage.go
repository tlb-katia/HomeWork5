package storage

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

func NewDB() (*sql.DB, error) {
	const op = "storage.NewDB"

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return db, nil
}

func CloseDB(db *sql.DB) error {
	const op = "storage.CloseDB"

	err := db.Close()
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	return nil
}
