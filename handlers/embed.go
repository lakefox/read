package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"main/models"
	"main/utils"
)

func EmbedAudio(w http.ResponseWriter, r *http.Request) {
	queryURL := r.URL.Query().Get("url")
	if queryURL == "" {
		queryURL = r.Referer()
		if queryURL == "" {
			http.Error(w, "URL not provided", http.StatusBadRequest)
			return
		}
	}

	db, err := utils.InitDB()
	if err != nil {
		log.Printf("Error initializing database: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var article models.Article
	var cachedURL string
	err = db.QueryRow("SELECT url, category, keywords, location, title, author, image, description, site FROM audio_cache WHERE url = ?", queryURL).Scan(
		&article.Url, &article.Category, &article.Keywords, &cachedURL, &article.Title, &article.Author, &article.Image, &article.Description, &article.Site)
	if err != nil {
		article, err = utils.FetchDoc(queryURL)
		if err != nil {
			log.Printf("Error fetching document: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		article.Url = queryURL
	}

	audioURL := fmt.Sprintf("/play?url=%s", queryURL)
	data := struct {
		AudioURL string
		Image    string
		Title    string
		Author   string
		Site     string
		Category string
		Date     string
	}{
		AudioURL: audioURL,
		Image:    article.Image,
		Title:    article.Title,
		Author:   article.Author,
		Site:     article.Site,
		Category: article.Category,
		Date:     article.Date,
	}

	tmpl, err := template.ParseFiles("/root/server/embed.html")
	if err != nil {
		log.Printf("Error parsing template: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
