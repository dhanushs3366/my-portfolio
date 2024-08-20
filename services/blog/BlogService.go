package blog

import (
	"database/sql"
	"dhanushs3366/my-portfolio/models"
	"log"
	"time"
)

type BlogStore struct {
	DB *sql.DB
}

func NewBlogStore(db *sql.DB) *BlogStore {
	return &BlogStore{
		DB: db,
	}
}

func (s *BlogStore) CreateBlogTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS BLOG(
			ID INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			USER_ID INT NOT NULL,
			IMAGES_KEY VARCHAR(255)[],
			CONTENT TEXT NOT NULL,
			CREATED_AT TIMESTAMP NOT NULL,
			UPDATED_AT TIMESTAMP NOT NULL,
			CONSTRAINT fk_user
				FOREIGN KEY (USER_ID)
				REFERENCES USERS(ID)
		)
	`

	_, err := s.DB.Exec(query)
	if err != nil {
		return err
	}

	log.Printf("BLOG table created")
	return nil
}

func (s *BlogStore) CreateBlog(user *models.User, content string) error {
	query := `
		INSERT INTO BLOG(USER_ID,IMAGES_KEY,CONTENT,CREATED_AT,UPDATED_AT)
		VALUES($1,$2,$3,$4,$5)
	`
	now := time.Now()
	_, err := s.DB.Exec(query, user.ID, nil, content, now, now)
	if err != nil {
		return err
	}
	log.Println("Blog created")
	return nil
}
