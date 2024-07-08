package handlers

import (
	"log"
	"net/http"

	"main/utils"
)

func StreamAudio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queryURL := r.URL.Query().Get("url")
	if queryURL == "" {
		queryURL = r.Referer()
		if queryURL == "" {
			http.Error(w, "URL not provided", http.StatusBadRequest)
			return
		}
	}

	log.Println(queryURL)

	db, err := utils.InitDB()
	if err != nil {
		log.Printf("Error initializing database: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var cachedURL string
	err = db.QueryRow("SELECT location FROM audio_cache WHERE url = ?", queryURL).Scan(&cachedURL)

	if err == nil {
		log.Printf("Serving File: %s", cachedURL)
		w.Header().Set("Content-Type", "audio/webm")
		http.ServeFile(w, r, cachedURL)
	} else {
		utils.StreamNewAudio(w, r, db, queryURL)
	}
}
