package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GithubRepo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	URL         string `json:"html_url"`
	Description string `json:"description"`
	CloneURL    string `json:"clone_url"`
	Language    string `json:"language"`
	Stars       uint   `json:"stargazers_count"`
	Watchers    uint   `json:"watchers_count"`
}

type GithubUser struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Bio       string `json:"bio"`
}

func FetchReposByUserName(username string) ([]GithubRepo, error) {
	URL, ok := os.LookupEnv("GITHUB_URL")
	if !ok {
		return nil, errors.New("env variable not found")
	}

	resp, err := http.Get(fmt.Sprintf("%s/users/%s/repos", URL, username))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var repos []GithubRepo

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
