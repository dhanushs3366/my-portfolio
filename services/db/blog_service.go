package db

import (
	"database/sql"
	"dhanushs3366/my-portfolio/models"
	"errors"
	"log"
	"time"
)

func (s *Store) CreateBlogTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS BLOG(
			ID INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			USER_ID INT NOT NULL,
			IMAGES_KEY VARCHAR(255)[],
			TITLE VARCHAR(255),
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

func (s *Store) CreateBlog(user *models.User, content string, title string) error {
	query := `
		INSERT INTO BLOG(USER_ID,IMAGES_KEY,CONTENT,CREATED_AT,UPDATED_AT,TITLE)
		VALUES($1,$2,$3,$4,$5,$6)
	`
	now := time.Now()
	_, err := s.DB.Exec(query, user.ID, nil, content, now, now, title)
	if err != nil {
		return err
	}
	log.Println("Blog created")
	return nil
}

func (s *Store) EditBlog(blogID string, content string, title string) error {
	query := `
		UPDATE BLOG
		SET CONTENT=$1,TITLE=$2,UPDATED_AT=$3,
		WHERE ID=$4
		AND DELETED IS NOT TRUE
	`

	_, err := s.DB.Exec(query, content, title, time.Now(), blogID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No rows found")
			return err
		}
		return err
	}
	return nil
}

func (s *Store) DeleteBlog(blogID string) error {
	query := `
		UPDATE BLOG 
		SET DELETED=TRUE,DELETED_AT=$1
		WHERE ID=$2
	`
	_, err := s.DB.Exec(query, time.Now(), blogID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No rows found")
			return err
		}
		return err
	}
	return nil
}

func (s *Store) GetBlogs() ([]models.Blog, error) {
	query := `
		SELECT ID,USER_ID,TITLE,CONTENT,CREATED_AT,UPDATED_AT FROM BLOG BL
		WHERE BL.DELETED=FALSE
	`

	var blogs []models.Blog

	rows, err := s.DB.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No rows found")
			return nil, err
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var blog models.Blog
		err = rows.Scan(&blog.ID, &blog.OwnedBy, &blog.Title, &blog.Content, &blog.CreatedAt, &blog.UpdatedAt)
		if err != nil {
			continue
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func (s *Store) GetBlogByID(ID string) (*models.Blog, error) {
	query := `
		SELECT ID,USER_ID,TITLE,CONTENT,CREATED_AT,UPDATED_AT FROM BLOG BL
		WHERE BL.DELETED=FALSE AND BL.ID=$1
	`
	row := s.DB.QueryRow(query, ID)
	var blog models.Blog
	err := row.Scan(&blog.ID, &blog.OwnedBy, &blog.Title, &blog.Content, &blog.CreatedAt, &blog.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No rows found")
			return nil, err
		}
		return nil, err
	}

	return &blog, nil
}
