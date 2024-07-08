package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/plugins/category"
	fixtext "main/plugins/date"
	"main/plugins/keywords"
	"main/plugins/rss"
	"net/http"
	"net/url"
	"strings"

	"main/models"
	"main/utils"
)

type RequestData struct {
	Data []string `json:"data"`
}

type ResponseData struct {
	Data []models.Article `json:"data"`
}

func Feed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method != http.MethodPost {
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

	var articles, oldArticles []models.Article
	for _, v := range urls {
		items, err := rss.ParseFeed(v)
		if err != nil {
			log.Fatalf("Error parsing RSS feed: %v", err)
		} else {
			for _, item := range items {
				parsedURL, err := url.Parse(item.Link)
				if err != nil {
					fmt.Println("Error parsing URL:", err)
					return
				}
				sp := strings.Split(parsedURL.Host, ".")
				hostname := sp[len(sp)-2]
				site := strings.ToUpper(hostname)

				text := item.Title + " " + item.Description
				doc := models.Article{
					Keywords:    keywords.ExtractKeywords(text, 3),
					Category:    category.GetCategory(text),
					Title:       item.Title,
					Author:      item.Author,
					Description: item.Description,
					Site:        site,
					Url:         item.Link,
					Date:        fixtext.FormatTime(&item.Published),
				}
				articles = append(articles, doc)
			}
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var requestData RequestData
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	for _, v := range requestData.Data {
		var article models.Article
		var cachedURL string
		err = db.QueryRow("SELECT url, category, keywords, location, title, author, image, description, site FROM audio_cache WHERE url = ?", v).Scan(
			&article.Url, &article.Category, &article.Keywords, &cachedURL, &article.Title, &article.Author, &article.Image, &article.Description, &article.Site)
		if err != nil {
			oldArticles = append(oldArticles, article)
		}
	}
	utils.ShuffleArticles(articles)
	utils.ShuffleArticles(oldArticles)

	sorted := utils.Rank(articles, oldArticles)
	responsePayload := ResponseData{Data: sorted[0:utils.Min(50, len(sorted))]}

	responseJSON, err := json.Marshal(responsePayload)
	if err != nil {
		http.Error(w, "Error creating response JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}