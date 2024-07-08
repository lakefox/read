package models

type Article struct {
	Keywords    []string `json:"keywords"`
	Voice       string   `json:"voice"`
	Category    string   `json:"category"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Image       string   `json:"image"`
	Description string   `json:"description"`
	Site        string   `json:"site"`
	Url         string   `json:"url"`
	Text        string   `json:"text"`
	Date        string   `json:"date"`
}
