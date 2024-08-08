package services

import (
	"database/sql"
	"dhanushs3366/my-portfolio/models"
	"errors"
	"log"
	"time"
)

var ErrNoEntityFound = errors.New("no entity found")

func createUserTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS USERS (
			ID SERIAL PRIMARY KEY,
			USERNAME VARCHAR(50) UNIQUE NOT NULL,
			PASSWORD VARCHAR(255) NOT NULL,
			IS_ADMIN BOOLEAN NOT NULL DEFAULT FALSE,
			CREATED_AT TIMESTAMP NOT NULL DEFAULT NOW(), 
			UPDATED_AT TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`

	_, err := DB.Exec(query)

	if err != nil {
		return err
	}

	log.Println("User table created")
	return nil
}

// before inserting user search if user already exists

func InsertUser(user *models.User) error {
	query := `
		INSERT INTO USERS (USERNAME,PASSWORD,IS_ADMIN,CREATED_AT,UPDATED_AT)
		VALUES ($1,$2,$3,$4,$5)
	`

	hasedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	now := time.Now()
	_, err = DB.Exec(query, user.Username, hasedPassword, user.IsAdmin, now, now)

	if err != nil {
		return err
	}

	log.Print("User created")
	return nil
}

func GetUser(username string) (*models.User, error) {
	query := `
		SELECT ID,USERNAME,PASSWORD,IS_ADMIN FROM USERS
		WHERE USERNAME=$1
	`
	row := DB.QueryRow(query, username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.IsAdmin)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoEntityFound
		}
		return nil, err
	}

	return &user, nil
}

func UpdatePassword(username string, password string) error {
	query := `
		UPDATE USERS SET PASSWORD=$1 WHERE USERNAME=$2
	`

	_, err := DB.Exec(query, password, username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoEntityFound
		}
		return err
	}

	return nil
}
