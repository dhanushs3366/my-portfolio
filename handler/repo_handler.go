package handler

import (
	"dhanushs3366/my-portfolio/models"
	"dhanushs3366/my-portfolio/services/hooks"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
)

func (h *Handler) fetchGitRepos(c echo.Context) error {
	repos, err := hooks.FetchGitRepos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, repos)
}

func (h *Handler) getValidRepos(c echo.Context) error {
	repos, err := h.store.GetValidRepos()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, repos)
}

func (h *Handler) updateRepo(c echo.Context) error {

	repoID, err := strconv.ParseUint(c.Param("repoID"), 10, 32)

	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid param repoID")
	}

	isVisible, err := strconv.ParseBool(c.QueryParam("isVisible"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid query param isVisible")
	}

	err = h.store.UpdateRepo(uint(repoID), isVisible)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "repo updated")
}

func (h *Handler) syncRepos(c echo.Context) error {
	prevRepos, err := h.store.GetAllRepos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to get previous repos from DB: %s", err.Error()))
	}

	repos, err := hooks.FetchGitRepos()
	if err != nil {
		return c.JSON(http.StatusFailedDependency, fmt.Sprintf("Failed to fetch repos from GitHub API: %s", err.Error()))
	}

	var updatedRepos []models.ValidRepo
	for index, repo := range repos {
		updatedRepos = append(updatedRepos, models.ValidRepo{
			ID:        index, // Placeholder
			RepoID:    repo.ID,
			IsVisible: false,
		})
	}

	newRepos := difference(updatedRepos, prevRepos)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for _, repo := range newRepos {
		wg.Add(1)
		go func(repoID uint) {
			defer wg.Done()
			if err := h.store.InsertRepos(repoID, false); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				log.Printf("Failed to insert repo: %s", err.Error())
			}
		}(repo.RepoID)
	}

	wg.Wait()

	if len(errs) > 0 {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Some repos failed to insert: %v", errs))
	}

	return c.JSON(http.StatusOK, newRepos)
}

// returns slice of a-b
// returns [] if a is a subset of b
func difference(a, b []models.ValidRepo) []models.ValidRepo {
	bMap := make(map[int]bool)
	for _, repo := range b {
		bMap[int(repo.RepoID)] = true
	}

	// Collect elements in a that are not in b
	var diff []models.ValidRepo
	for _, repo := range a {
		if !bMap[int(repo.RepoID)] {
			diff = append(diff, repo)
		}
	}
	return diff
}
