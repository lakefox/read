package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"main/utils"
)

func AddRSS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	db, err := utils.InitDB()
	if err != nil {
		log.Printf("Error initializing database: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	queryURL := r.URL.Query().Get("url")
	if queryURL == "" {
		http.Error(w, "URL not provided", http.StatusBadRequest)
		return
	}

	log.Println(queryURL)

	rows, err := db.Query("SELECT url FROM rss_feeds")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			log.Fatal(err)
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	var exists bool
	for _, v := range urls {
		if v == queryURL {
			exists = true
			break
		}
	}
	if !exists {
		_, err = db.Exec("INSERT INTO rss_feeds (url) VALUES (?)", queryURL)
		if err != nil {
			log.Printf("Error inserting into database: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		urls = append(urls, queryURL)
	}

	responseJSON, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, "Error creating response JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
