package blog

import (
	"database/sql"
	"dhanushs3366/my-portfolio/models"
	"dhanushs3366/my-portfolio/services/db"
	"errors"
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
			DELETED BOOLEAN DEFAULT FALSE,
			DELETED_AT TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_blog_user
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

func (s *BlogStore) EditBlog(blogID string, content string) error {
	query := `
		UPDATE BLOG
		SET CONTENT=$1,UPDATED_AT=$2
		WHERE ID=$3
		AND DELETED IS NOT TRUE
	`

	_, err := s.DB.Exec(query, content, time.Now(), blogID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.ErrNoEntityFound
		}
		return err
	}
	return nil
}

func (s *BlogStore) DeleteBlog(blogID string) error {
	query := `
		UPDATE BLOG 
		SET DELETED=TRUE,DELETED_AT=$1
		WHERE ID=$2
	`
	_, err := s.DB.Exec(query, time.Now(), blogID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.ErrNoEntityFound
		}
		return err
	}
	return nil
}

func (s *BlogStore) GetBlogs() ([]models.Blog, error) {
	query := `
		SELECT ID,USER_ID,CONTENT,CREATED_AT,UPDATED_AT FROM BLOG BL
		WHERE BL.DELETED=FALSE
	`

	var blogs []models.Blog

	rows, err := s.DB.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, db.ErrNoEntityFound
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var blog models.Blog
		err = rows.Scan(&blog.ID, &blog.OwnedBy, &blog.Content, &blog.CreatedAt, &blog.UpdatedAt)
		if err != nil {
			continue
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}
