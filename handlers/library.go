package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"main/models"
	"main/utils"
)

func HandleLibrary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) != 2 || pathParts[0] != "library" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	libraryName := pathParts[1]

	db, err := utils.InitDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodPost:
		addArticleToLibrary(w, r, db, libraryName)
	case http.MethodGet:
		getLibraryArticles(w, r, db, libraryName)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func addArticleToLibrary(w http.ResponseWriter, r *http.Request, db *sql.DB, libraryName string) {
	var requestBody struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert the URL into the library table
	insertQuery := `INSERT INTO library (library_name, url) VALUES (?, ?)`
	_, err = db.Exec(insertQuery, libraryName, requestBody.URL)
	if err != nil {
		http.Error(w, "Error adding article to library", http.StatusInternalServerError)
		return
	}

	// Start the audio download in a new goroutine
	go func() {
		err := utils.DownloadAudio(db, requestBody.URL)
		if err != nil {
			log.Printf("Error downloading audio: %v\n", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Article added to library and audio download started")
}

func getLibraryArticles(w http.ResponseWriter, r *http.Request, db *sql.DB, libraryName string) {
	query := `
	SELECT a.url, a.category, a.keywords, a.location, a.title, a.author, a.image, a.description, a.site
	FROM library l
	JOIN audio_cache a ON l.url = a.url
	WHERE l.library_name = ?
	`
	rows, err := db.Query(query, libraryName)
	if err != nil {
		http.Error(w, "Error querying library articles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		var keywords string
		if err := rows.Scan(&article.Url, &article.Category, &keywords, &article.Url, &article.Title, &article.Author, &article.Image, &article.Description, &article.Site); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		json.Unmarshal([]byte(keywords), &article.Keywords)
		articles = append(articles, article)
	}

	responseJSON, err := json.Marshal(articles)
	if err != nil {
		http.Error(w, "Error creating response JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
