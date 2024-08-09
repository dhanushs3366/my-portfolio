package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var ErrNoEntityFound = errors.New("no entity found in the DB")
var db *sql.DB

func ConnectToDB() (*sql.DB, error) {
	userName := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	DBName := os.Getenv("DB_NAME")
	DBHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", userName, password, DBName, DBHost, DBPort)

	var err error
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	err = db.Ping()

	if err != nil {
		return nil, err
	}

	log.Println("Connected to DB sucessfully")

	return db, nil
}

func Close() error {
	return db.Close()
}
