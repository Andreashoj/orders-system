package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
	query, err := os.ReadFile("internal/db/init.sql")
	if err != nil {
		return nil, fmt.Errorf("failed reading sql from init.sql: %s", err)
	}

	conn := "postgres://postgres:postgres@localhost:5432/orders_db?sslmode=disable"
	DB, err := sql.Open("postgres", conn)

	if err != nil {
		return nil, fmt.Errorf("failed starting DB client: %s", err)
	}

	_, err = DB.Exec(string(query))
	if err != nil {
		return nil, fmt.Errorf("failed seeding database with init.sql seed file: %s", err)
	}

	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetConnMaxIdleTime(5 * time.Minute)

	return DB, err
}
