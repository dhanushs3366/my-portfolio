package user

import (
	"database/sql"
	"dhanushs3366/my-portfolio/models"
	"errors"
	"log"
	"time"
)

var ErrNoEntityFound = errors.New("no entity found in DB")

type UserStore struct {
	DB *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		DB: db,
	}
}

// uni placements are happening i need to learn sql so im raw dogging this without any orms

// create a new table entry once in a while like every 1 hour update the latest record for the remaining time

func (s *UserStore) CreateUserTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS USERS(
			ID SERIAL PRIMARY KEY,
			USERNAME VARCHAR(50) UNIQUE NOT NULL,
			PASSWORD VARCHAR(255) NOT NULL,
			IS_ADMIN BOOLEAN NOT NULL DEFAULT FALSE,
			CREATED_AT TIMESTAMP NOT NULL DEFAULT NOW(),
			UPDATED_AT TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
	_, err := s.DB.Exec(query)

	if err != nil {
		return err
	}

	log.Println("user table created ")
	return nil
}

func (s *UserStore) InsertUser(username string, password string, isAdmin bool) error {
	query := `
		INSERT INTO USERS (USERNAME,PASSWORD,IS_ADMIN,CREATED_AT,UPDATED_AT)
		VALUES ($1,$2,$3,$4,$5)
	`

	now := time.Now()
	_, err := s.DB.Exec(query, username, password, isAdmin, now, now)

	if err != nil {
		return err
	}

	log.Print("User created")
	return nil
}

func (s *UserStore) GetUser(username string) (*models.User, error) {
	query := `
		SELECT ID,USERNAME,PASSWORD,IS_ADMIN FROM USERS
		WHERE USERNAME=$1
	`
	row := s.DB.QueryRow(query, username)
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

func (s *UserStore) UpdatePassword(username string, password string) error {
	query := `
		UPDATE USERS SET PASSWORD=$1, UPDATED_AT=$2 WHERE USERNAME=$3
	`

	_, err := s.DB.Exec(query, password, time.Now(), username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoEntityFound
		}
		return err
	}

	return nil
}
