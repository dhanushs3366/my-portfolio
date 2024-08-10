package logger

import (
	"database/sql"
	"dhanushs3366/my-portfolio/models"
	"errors"
	"log"
	"time"
)

type LogStore struct {
	DB *sql.DB
}

func NewLogStore(db *sql.DB) *LogStore {
	return &LogStore{
		DB: db,
	}
}

func (s *LogStore) CreateLogActivityTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS LOG_ACTIVITY(
			ID SERIAL PRIMARY KEY,
			KEY INTEGER,
			MIDDLE_CLICKS INTEGER,
			RIGHT_CLICKS INTEGER,
			LEFT_CLICKS INTEGER,
			EXTRA_CLICKS INTEGER,
			CREATED_AT TIMESTAMP NOT NULL,
			UPDATED_AT TIMESTAMP NOT NULL
		)
	`

	_, err := s.DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Log Activity table created")
	return nil
}

func (s *LogStore) InsertLogActivity(activity *models.LoggedActivity) error {
	query := `
		INSERT INTO LOG_ACTIVITY (key,middle_clicks,right_clicks,left_clicks,extra_clicks,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`
	now := time.Now()
	result, err := s.DB.Exec(query, activity.Key, activity.MiddleClicks, activity.RightClicks, activity.LeftClicks, activity.ExtraClicks, now, now)
	if err != nil {
		return err
	}
	lastRow, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Printf("Inserted into log activity, last inserted row: %d", lastRow)

	return nil
}

func (s *LogStore) GetRecentUpdatedActivity() (*models.LoggedActivity, error) {
	query := `
		SELECT key, middle_clicks, right_clicks,left_clicks,extra_clicks FROM LOG_ACTIVITY
		ORDER BY created_at DESC
		LIMIT 1
	`

	row := s.DB.QueryRow(query)

	var activity models.LoggedActivity

	err := row.Scan(&activity.Key, &activity.RightClicks, &activity.LeftClicks, &activity.MiddleClicks, &activity.ExtraClicks)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &activity, nil
}

func (s *LogStore) GetLogActivityById(ID int) (*models.LoggedActivity, error) {
	query := `
		select key, middle_clicks,left_clicks,right_clicks,extra_clicks, created_at,updated_at FROM LOG_ACTIVITY
		WHERE id=$1
	`
	row := s.DB.QueryRow(query, ID)
	var activity models.LoggedActivity
	err := row.Scan(&activity.Key, &activity.MiddleClicks, &activity.LeftClicks, &activity.RightClicks, &activity.ExtraClicks, &activity.CreatedAt, &activity.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no rows in the table")

		}
		return nil, err
	}
	return &activity, nil
}

func (s *LogStore) UpdateLogActivityById(ID int, updatedActivity models.LoggedActivity) error {
	activity, err := s.GetLogActivityById(ID)
	if err != nil {
		return err
	}

	updatedActivity.Key += activity.Key
	updatedActivity.LeftClicks += activity.LeftClicks
	updatedActivity.RightClicks += activity.RightClicks
	updatedActivity.MiddleClicks += activity.MiddleClicks
	updatedActivity.ExtraClicks += activity.ExtraClicks

	query := `
		UPDATE LOG_ACTIVITY
		SET key=$1, left_clicks=$2, right_clicks=$3,middle_clicks=$4,extra_clicks=$5, updated_at=$6
		WHERE ID=$7
	`

	_, err = s.DB.Exec(query, updatedActivity.Key, updatedActivity.LeftClicks, updatedActivity.RightClicks, updatedActivity.MiddleClicks, updatedActivity.ExtraClicks, time.Now(), ID)
	if err != nil {
		return err
	}

	log.Printf("Updated the record with id %d\n", ID)
	return nil

}

func (s *LogStore) GetRecentLogActivityCreatedAt() (int, *time.Time, error) {
	query := `
		SELECT id, created_at FROM LOG_ACTIVITY
		ORDER BY created_at DESC
		LIMIT 1 
	`

	row := s.DB.QueryRow(query)
	var id int
	var createdAt time.Time
	err := row.Scan(&id, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil, nil
		}
		return 0, nil, err
	}

	return id, &createdAt, nil
}

func (s *LogStore) GetLogActivtyPerWeek(toDate time.Time) ([]models.LoggedActivity, error) {
	var logs []models.LoggedActivity
	query := `
		SELECT * FROM LOG_ACTIVITY 
		WHERE UPDATED_AT >=$1 AND UPDATED_AT <=$2
	`
	fromWeek := toDate.Add(-7 * 24 * time.Hour)
	rows, err := s.DB.Query(query, fromWeek, toDate)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var log models.LoggedActivity
		err := rows.Scan(
			&log.ID,
			&log.Key,
			&log.MiddleClicks,
			&log.RightClicks,
			&log.LeftClicks,
			&log.ExtraClicks,
			&log.CreatedAt,
			&log.UpdatedAt,
		)

		if err != nil {
			continue
		}
		logs = append(logs, log)
	}
	if err = rows.Err(); err != nil {
		return nil, err

	}

	return logs, nil
}
