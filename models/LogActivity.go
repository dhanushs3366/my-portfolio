package models

import (
	"log"

	_ "github.com/lib/pq"
)

// uni placements are happening i need to learn sql so im raw dogging this without any orms

/*

type Logger struct {
	Key          int `json:"all_keys"`
	MiddleClicks int `json:"middle_clicks"`
	RightClicks  int `json:"right_clicks"`
	LeftClicks   int `json:"left_clicks"`
	ExtraClicks  int `json:"extra_clicks"`
}
*/
// create a new table entry once in a while like every 1 hour update the latest record for the remaining time

func createTempTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS temp(
        id SERIAL PRIMARY KEY
    )
    `
	record, err := DB.Exec(query)
	if err != nil {
		return err
	}

	log.Printf("Tablel created %+v\n", record)
	return nil
}
