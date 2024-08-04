package models

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// uni placements are happening i need to learn sql so im raw dogging this without any orms

type LoggedActivity struct {
	Key          int `json:"all_keys"`
	MiddleClicks int `json:"middle_clicks"`
	RightClicks  int `json:"right_clicks"`
	LeftClicks   int `json:"left_clicks"`
	ExtraClicks  int `json:"extra_clicks"`
}

// create a new table entry once in a while like every 1 hour update the latest record for the remaining time

func createLogActivityTable() error {

	query := `
		CREATE TABLE IF NOT EXISTS LOG_ACTIVITY(
			id SERIAL PRIMARY KEY,
			key INTEGER,
			middle_clicks INTEGER,
			right_clicks INTEGER,
			left_clicks INTEGER,
			extra_clicks INTEGER,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`

	result, err := DB.Exec(query)
	if err != nil {
		return err
	}
	log.Printf("Log Activity table created %v\n", result)
	return nil
}

func InsertLogActivity(activity *LoggedActivity) error {
	query := `
		INSERT INTO LOG_ACTIVITY (key,middle_clicks,right_clicks,left_clicks,extra_clicks,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`
	now := time.Now()
	result, err := DB.Exec(query, activity.Key, activity.MiddleClicks, activity.RightClicks, activity.LeftClicks, activity.ExtraClicks, now, now)
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

func GetRecentUpdatedActivity() (*LoggedActivity, error) {
	query := `
		SELECT key, middle_clicks, right_clicks,left_clicks,extra_clicks FROM LOG_ACTIVITY
		ORDER BY created_at DESC
		LIMIT 1
	`

	row := DB.QueryRow(query)

	var activity LoggedActivity

	err := row.Scan(&activity.Key, &activity.RightClicks, &activity.LeftClicks, &activity.MiddleClicks, &activity.ExtraClicks)

	if err != nil {
		if err == sql.ErrNoRows {
			// -1 with no error returned indicates no rows
			return nil, errors.New("no rows in the table")
		}
		return nil, err
	}

	return &activity, nil
}

func GetLogActivityById(ID int) (*LoggedActivity, error) {
	query := `
		select key, middle_clicks,left_clicks,right_clicks,extra_clicks FROM LOG_ACTIVITY
		WHERE id=$1
	`
	row := DB.QueryRow(query, ID)
	var activity LoggedActivity
	err := row.Scan(&activity.Key, &activity.MiddleClicks, &activity.LeftClicks, &activity.RightClicks, &activity.ExtraClicks)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no rows in the table")

		}
		return nil, err
	}
	return &activity, nil
}

func UpdateLogActivityById(ID int, updatedActivity LoggedActivity) error {
	activity, err := GetLogActivityById(ID)
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

	_, err = DB.Exec(query, updatedActivity.Key, updatedActivity.LeftClicks, updatedActivity.RightClicks, updatedActivity.MiddleClicks, updatedActivity.MiddleClicks, time.Now(), ID)
	if err != nil {
		return err
	}

	log.Printf("Updated the record with id %d\n", ID)
	return nil

}

func GetRecentLogActivityCreatedAt() (int, *time.Time, error) {
	query := `
		SELECT id, created_at FROM LOG_ACTIVITY
		ORDER BY created_at DESC
		LIMIT 1 
	`

	row := DB.QueryRow(query)

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
