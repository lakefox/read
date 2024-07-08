package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"main/models"
	"main/utils"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"main/plugins/category"
	fixtext "main/plugins/date"
	"main/plugins/keywords"
	"main/plugins/rss"
)

type RequestData struct {
	Data []string `json:"data"`
}

type ResponseData struct {
	Data []models.Article `json:"data"`
}

func Feed(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

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
		for _, feedURL := range urls {
			cacheFile := getCacheFileName(feedURL)
			if fileExists(cacheFile) && !isFileExpired(cacheFile, 6*time.Hour) {
				cachedArticles, err := readArticlesFromFile(cacheFile)
				if err == nil {
					articles = append(articles, cachedArticles...)
					continue
				}
			}

			items, err := rss.ParseFeed(feedURL)
			if err != nil {
				log.Fatalf("Error parsing RSS feed: %v", err)
			} else {
				var feedArticles []models.Article
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
					feedArticles = append(feedArticles, doc)
				}
				articles = append(articles, feedArticles...)
				saveArticlesToFile(cacheFile, feedArticles)
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

		var sorted []models.Article

		if len(requestData.Data) > 0 {
			for _, v := range requestData.Data {
				var article models.Article
				var cachedURL string
				err = db.QueryRow("SELECT url, category, keywords, location, title, author, image, description, site FROM audio_cache WHERE url = ?", v).Scan(
					&article.Url, &article.Category, &article.Keywords, &cachedURL, &article.Title, &article.Author, &article.Image, &article.Description, &article.Site)
				if err != nil {
					oldArticles = append(oldArticles, article)
				}
			}
			utils.ShuffleArticles(oldArticles)
			utils.ShuffleArticles(articles)
			sorted = utils.Rank(articles, oldArticles)
		} else {
			utils.ShuffleArticles(articles)
			sorted = articles
		}

		responsePayload := ResponseData{Data: sorted[0:utils.Min(50, len(sorted))]}

		responseJSON, err := json.Marshal(responsePayload)
		if err != nil {
			http.Error(w, "Error creating response JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}
}

func getCacheFileName(url string) string {
	return filepath.Join("/tmp", fmt.Sprintf("%x.json", url))
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func isFileExpired(filename string, maxAge time.Duration) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return true
	}
	return time.Since(info.ModTime()) > maxAge
}

func readArticlesFromFile(filename string) ([]models.Article, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var articles []models.Article
	if err := json.Unmarshal(data, &articles); err != nil {
		return nil, err
	}
	return articles, nil
}

func saveArticlesToFile(filename string, articles []models.Article) error {
	data, err := json.Marshal(articles)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}
