package handler

import (
	"dhanushs3366/my-portfolio/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handler) PostLogDetails(c echo.Context) error {
	var loggedActivity models.LoggedActivity
	err := json.NewDecoder(c.Request().Body).Decode(&loggedActivity)
	if err != nil {
		return err
	}
	log.Print(loggedActivity)

	lastID, lastCreatedAt, err := h.logStore.GetRecentLogActivityCreatedAt()
	if err != nil {
		return err
	}

	// -1 id means no rows in the table
	if lastID == -1 || time.Since(*lastCreatedAt) > time.Minute {
		log.Printf("last updated row was 1 hr ago or no rows exist, creating a new row")
		err := h.logStore.InsertLogActivity(&loggedActivity)
		if err != nil {
			return err
		}
	} else {
		// update the latest record in the table
		err = h.logStore.UpdateLogActivityById(lastID, loggedActivity)
		if err != nil {
			return err
		}
		log.Printf("Updated ID:%d record", lastID)
	}

	return nil
}

func (h *Handler) GetLogDetails(c echo.Context) error {
	to := c.QueryParam("to")
	dateFormat := "02-01-2006"

	if to == "" {
		return c.JSON(http.StatusNoContent, "query param not found")
	}

	toDate, err := time.Parse(dateFormat, to)

	if err != nil {
		return c.JSON(http.StatusConflict, err.Error())
	}

	logs, err := h.logStore.GetLogActivtyPerWeek(toDate)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, logs)
}
