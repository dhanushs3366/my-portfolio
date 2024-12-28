package models

// github repo is not saved in the db
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

type ValidRepo struct {
	ID        int  `json:"id"`
	RepoID    uint `json:"repo_id"`
	IsVisible bool `json:"is_visible"`
}
