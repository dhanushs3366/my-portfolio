package api

import (
	"dhanushs3366/my-portfolio/models"
	"encoding/json"
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func postLogDetails(c echo.Context) error {
	var loggedActivity models.LoggedActivity
	err := json.NewDecoder(c.Request().Body).Decode(&loggedActivity)
	if err != nil {
		return err
	}
	log.Print(loggedActivity)

	lastID, lastCreatedAt, err := models.GetRecentLogActivityCreatedAt()
	if err != nil {
		return err
	}

	// -1 id means no rows in the table
	if lastCreatedAt.Compare(time.Now()) > int(time.Hour) || lastID == -1 {
		log.Printf("last updated row was 1 hr ago creating a new row")
		err := models.InsertLogActivity(&loggedActivity)
		if err != nil {
			return err
		}
	}

	// update the latest record in the table
	err = models.UpdateLogActivityById(lastID, loggedActivity)
	if err != nil {
		return err
	}
	log.Printf("Updated ID:%d record", lastID)
	return nil
}
