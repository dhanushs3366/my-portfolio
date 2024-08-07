package handlers

import (
	"dhanushs3366/my-portfolio/models"
	"dhanushs3366/my-portfolio/services"
	"encoding/json"
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func PostLogDetails(c echo.Context) error {
	var loggedActivity models.LoggedActivity
	err := json.NewDecoder(c.Request().Body).Decode(&loggedActivity)
	if err != nil {
		return err
	}
	log.Print(loggedActivity)

	lastID, lastCreatedAt, err := services.GetRecentLogActivityCreatedAt()
	if err != nil {
		return err
	}

	// -1 id means no rows in the table
	if lastID == -1 || time.Since(*lastCreatedAt) > time.Minute {
		log.Printf("last updated row was 1 hr ago or no rows exist, creating a new row")
		err := services.InsertLogActivity(&loggedActivity)
		if err != nil {
			return err
		}
	} else {
		// update the latest record in the table
		err = services.UpdateLogActivityById(lastID, loggedActivity)
		if err != nil {
			return err
		}
		log.Printf("Updated ID:%d record", lastID)
	}

	return nil
}
