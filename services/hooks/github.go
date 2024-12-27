package hooks

import (
	"dhanushs3366/my-portfolio/api"
	"dhanushs3366/my-portfolio/models"
	"log"
	"os"
)

// idk why i made a sub packages for db,users,hooks,blog etc but im too lazy to refactor it
// github token expires so i will just hardcode my github username
func FetchGitRepos() ([]models.GithubRepo, error) {
	username := os.Getenv("GITHUB_USERNAME")

	repos, err := api.FetchReposByUserName(username)

	if err != nil {
		log.Printf("Error fetching user repos")
		return nil, err
	}

	return repos, nil
}
