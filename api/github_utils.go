package api

import (
	"dhanushs3366/my-portfolio/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GithubUser struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Bio       string `json:"bio"`
}

func FetchReposByUserName(username string) ([]models.GithubRepo, error) {
	URL, ok := os.LookupEnv("GITHUB_URL")
	if !ok {
		return nil, errors.New("env variable not found")
	}

	resp, err := http.Get(fmt.Sprintf("%s/users/%s/repos", URL, username))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var repos []models.GithubRepo

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &repos)

	if err != nil {
		return nil, err
	}

	return repos, err
}

// get user by github accesstoken instead of hardcoding my username
// future proof if i ever change my user name :)

func FetchUser() (*GithubUser, error) {
	URL := os.Getenv("GITHUB_URL")
	TOKEN := os.Getenv("GITHUB_ACCESS_TOKEN")

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", TOKEN))
	req.Header.Set("X-Github-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user GithubUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
