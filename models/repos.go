package models

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
	Repo      GithubRepo `json:"repo"`
	IsVisible bool       `json:"is_visible"`
}
