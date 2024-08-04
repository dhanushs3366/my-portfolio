package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var DB *sql.DB

func connectToDB() error {
	userName := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	DBName := os.Getenv("DB_NAME")
	DBHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", userName, password, DBName, DBHost, DBPort)
	log.Println("Connecting to DB")

	var err error
	DB, err = sql.Open("postgres", connStr)
	// do not use := to initialise DB it makes the scope of the var local you cant acess it in other places
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Println("Connected to DB successfully")
	return nil
}
func Init() error {
	err := connectToDB()
	// add a func call to sync database
	if err != nil {
		return err
	}
	err = createTempTable()
	if err != nil {
		return err
	}
	return nil
}
