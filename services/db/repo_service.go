package db

import (
	"database/sql"
	"dhanushs3366/my-portfolio/models"
	"errors"
	"log"
)

func (s *Store) createReposTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS REPOS(
			ID INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			REPO_ID INT UNIQUE NOT NULL,
			IS_VISIBLE BOOLEAN NOT NULL DEFAULT FALSE
		)
	`

	_, err := s.DB.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

// need an insert and update and prolly delete for valid repos
func (s *Store) InsertRepos(repoID uint, isVisible bool) error {
	query := `
		INSERT INTO REPOS(REPO_ID, IS_VISIBLE)
		VALUES($1,$2)
	`
	_, err := s.DB.Exec(query, repoID, isVisible)

	if err != nil {
		log.Printf("couldnt insert repo")
		return err
	}

	return nil
}

func (s *Store) GetRepo(repoID uint) (*models.ValidRepo, error) {
	query := `
		SELECT * FROM REPOS
		WHERE REPO_ID=$1
	`
	row := s.DB.QueryRow(query, repoID)
	var repo models.ValidRepo

	err := row.Scan(&repo.ID, &repo.RepoID, &repo.IsVisible)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No rows found")
		}
		return nil, err
	}

	return &repo, err
}

func (s *Store) GetValidRepos() ([]models.ValidRepo, error) {
	query := `
		SELECT * FROM REPOS
		WHERE IS_VISIBLE=TRUE
	`
	rows, err := s.DB.Query(query)

	if err != nil {
		return nil, err
	}

	var repos []models.ValidRepo

	defer rows.Close()

	for rows.Next() {
		var repo models.ValidRepo

		err = rows.Scan(&repo.ID, &repo.RepoID, &repo.IsVisible)

		if err != nil {
			log.Printf("Cant scan row")
			continue
		}

		repos = append(repos, repo)
	}

	return repos, nil
}

// only update you can do for this change visibility
func (s *Store) UpdateRepo(repoID uint, isVisible bool) error {
	query := `
		UPDATE REPOS
		SET IS_VISIBLE=$1
		WHERE REPO_ID=$2
	`
	_, err := s.DB.Exec(query, isVisible, repoID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("no rows found")
		}
		return err
	}

	return nil
}

func (s *Store) GetAllRepos() ([]models.ValidRepo, error) {
	query := `
		SELECT * FROM REPOS
	`

	rows, err := s.DB.Query(query)

	if err != nil {
		return nil, err
	}

	var repos []models.ValidRepo

	for rows.Next() {
		var repo models.ValidRepo

		err = rows.Scan(&repo.ID, &repo.RepoID, &repo.IsVisible)

		if err != nil {
			log.Printf("Cant scan row")
			continue
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

// no delete function is required as i would want all my git repos visibility to be stored in db
