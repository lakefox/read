package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"main/plugins/category"
	fixtext "main/plugins/date"
	"main/plugins/gender"
	"main/plugins/keywords"
	"main/plugins/rss"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	readability "github.com/go-shiori/go-readability"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

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

func fetchDoc(url string) (Article, error) {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return Article{}, fmt.Errorf("failed to parse %s, %v", url, err)
	}
	text := article.TextContent
	profile := gender.Profile(text)
	return Article{
		Keywords:    keywords.ExtractKeywords(text, 3),
		Voice:       profile["gender"],
		Category:    category.GetCategory(text),
		Title:       article.Title,
		Author:      article.Byline,
		Image:       article.Image,
		Description: article.Excerpt,
		Site:        article.SiteName,
		Url:         url,
		Date:        fixtext.FormatTime(article.PublishedTime),
		Text:        fixtext.ReplaceDates(text),
	}, nil
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "/root/server/audio.db")
	if err != nil {
		return nil, err
	}
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS audio_cache (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		category TEXT NOT NULL,
		keywords TEXT, -- Store keywords as a JSON string
		location TEXT NOT NULL,
		title TEXT,
		author TEXT,
		image TEXT,
		description TEXT,
		site TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS rss_feeds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func loadEnv() error {
	err := godotenv.Load("/root/server/.env")
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}

func streamAudio(w http.ResponseWriter, r *http.Request) {
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

	db, err := initDB()
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
		http.ServeFile(w, r, cachedURL)
	} else {
		opt, err := fetchDoc(queryURL)
		if err != nil {
			log.Printf("Error fetching document: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var voiceModel string
		if opt.Voice == "male" {
			voiceModel = "en_US-joe-medium.onnx"
		} else {
			voiceModel = "en_US-amy-medium.onnx"
		}
		log.Printf("Using voice model: %s\n", voiceModel)

		rawFile := fmt.Sprintf("/tmp/audio_%d.raw", time.Now().UnixNano())
		webmFile := fmt.Sprintf("/tmp/audio_%d.webm", time.Now().UnixNano())
		defer os.Remove(rawFile)
		// defer os.Remove(webmFile)

		piperCmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' | /root/piper/piper --model /root/%s --output-raw", opt.Text, voiceModel))
		ffmpegCmd := exec.Command("ffmpeg", "-f", "s16le", "-ar", "22050", "-ac", "1", "-i", "pipe:0", "-c:a", "libopus", "-f", "webm", "pipe:1")

		piperStdout, err := piperCmd.StdoutPipe()
		if err != nil {
			log.Printf("Error getting stdout pipe: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		ffmpegCmd.Stdin = piperStdout

		ffmpegStdout, err := ffmpegCmd.StdoutPipe()
		if err != nil {
			log.Printf("Error getting stdout pipe: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "audio/webm")
		w.Header().Set("Transfer-Encoding", "chunked")
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		if err := piperCmd.Start(); err != nil {
			log.Printf("Error starting Piper command: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := ffmpegCmd.Start(); err != nil {
			log.Printf("Error starting FFmpeg command: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		go func() {
			outFile, err := os.Create(webmFile)
			if err != nil {
				log.Printf("Error creating output file: %v\n", err)
				return
			}
			defer outFile.Close()

			multiWriter := io.MultiWriter(outFile, w)
			_, err = io.Copy(multiWriter, ffmpegStdout)
			if err != nil {
				log.Printf("Error copying data: %v\n", err)
			}
			flusher.Flush()
		}()

		if err := piperCmd.Wait(); err != nil {
			log.Printf("Error waiting for Piper command: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := ffmpegCmd.Wait(); err != nil {
			log.Printf("Error waiting for FFmpeg command: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		keywordsJSON, err := json.Marshal(opt.Keywords)
		if err != nil {
			log.Printf("Error marshaling keywords: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO audio_cache (url, category, keywords, location, title, author, image, description, site) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			queryURL, opt.Category, string(keywordsJSON), webmFile, opt.Title, opt.Author, opt.Image, opt.Description, opt.Site)
		if err != nil {
			log.Printf("Error inserting into database: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func embedAudio(w http.ResponseWriter, r *http.Request) {
	queryURL := r.URL.Query().Get("url")
	if queryURL == "" {
		queryURL = r.Referer()
		if queryURL == "" {
			http.Error(w, "URL not provided", http.StatusBadRequest)
			return
		}
	}

	db, err := initDB()
	if err != nil {
		log.Printf("Error initializing database: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var article Article
	var cachedURL string
	err = db.QueryRow("SELECT url, category, keywords, location, title, author, image, description, site FROM audio_cache WHERE url = ?", queryURL).Scan(
		&article.Url, &article.Category, &article.Keywords, &cachedURL, &article.Title, &article.Author, &article.Image, &article.Description, &article.Site)
	if err != nil {
		article, err = fetchDoc(queryURL)
		if err != nil {
			log.Printf("Error fetching document: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		article.Url = queryURL
	}

	audioURL := fmt.Sprintf("/?url=%s", queryURL)
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

type RequestData struct {
	Data []string `json:"data"`
}

type ResponseData struct {
	Data []Article `json:"data"`
}

func feed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	db, err := initDB()
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

	// Slice to hold the URLs
	var urls []string

	// Iterate over the rows and append URLs to the slice
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			log.Fatal(err)
		}
		urls = append(urls, url)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	var articles []Article
	var oldArticles []Article
	for _, v := range urls {
		items, err := rss.ParseFeed(v)
		if err != nil {
			log.Fatalf("Error parsing RSS feed: %v", err)
		} else {
			// item.Title, item.Link, item.Description, item.Published, item.Author
			for _, item := range items {
				parsedURL, err := url.Parse(item.Link)
				if err != nil {
					fmt.Println("Error parsing URL:", err)
					return
				}
				// Extract the hostname
				sp := strings.Split(parsedURL.Host, ".")
				hostname := sp[len(sp)-2]

				// Capitalize the hostname
				site := strings.ToUpper(hostname)

				text := item.Title + " " + item.Description

				doc := Article{
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

	// Now we have the feeds

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the JSON body into the RequestData struct
	var requestData RequestData
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	for _, v := range requestData.Data {
		// Get the info of the posts in the cache
		var article Article
		var cachedURL string
		err = db.QueryRow("SELECT url, category, keywords, location, title, author, image, description, site FROM audio_cache WHERE url = ?", v).Scan(
			&article.Url, &article.Category, &article.Keywords, &cachedURL, &article.Title, &article.Author, &article.Image, &article.Description, &article.Site)
		if err != nil {
			oldArticles = append(oldArticles, article)
		}
	}

	sorted := rank(articles, oldArticles)
	// Create the response payload
	responsePayload := ResponseData{Data: sorted}

	// Convert the response payload to JSON
	responseJSON, err := json.Marshal(responsePayload)
	if err != nil {
		http.Error(w, "Error creating response JSON", http.StatusInternalServerError)
		return
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON response
	w.Write(responseJSON)
}

func rank(newArticles, oldArticles []Article) []Article {
	// Filter out articles with duplicate URLs
	filteredArticles := filterDuplicateUrls(newArticles, oldArticles)

	// Calculate similarity scores
	type ScoredArticle struct {
		Article Article
		Score   int
	}
	scoredArticles := []ScoredArticle{}
	for _, newArticle := range filteredArticles {
		score := calculateMaxSimilarity(newArticle, oldArticles)
		scoredArticles = append(scoredArticles, ScoredArticle{Article: newArticle, Score: score})
	}

	// Sort articles by score in descending order
	sort.Slice(scoredArticles, func(i, j int) bool {
		return scoredArticles[i].Score > scoredArticles[j].Score
	})

	// Extract sorted articles
	sortedArticles := []Article{}
	for _, scoredArticle := range scoredArticles {
		sortedArticles = append(sortedArticles, scoredArticle.Article)
	}

	return sortedArticles
}

func filterDuplicateUrls(newArticles, oldArticles []Article) []Article {
	oldUrls := make(map[string]bool)
	for _, article := range oldArticles {
		oldUrls[article.Url] = true
	}

	filteredArticles := []Article{}
	for _, article := range newArticles {
		if !oldUrls[article.Url] {
			filteredArticles = append(filteredArticles, article)
		}
	}

	return filteredArticles
}

func calculateMaxSimilarity(newArticle Article, oldArticles []Article) int {
	maxScore := 0
	for _, oldArticle := range oldArticles {
		score := calculateSimilarity(newArticle, oldArticle)
		if score > maxScore {
			maxScore = score
		}
	}
	return maxScore
}

func calculateSimilarity(a, b Article) int {
	score := 0

	// Compare Keywords
	keywordSet := make(map[string]bool)
	for _, keyword := range a.Keywords {
		keywordSet[keyword] = true
	}
	for _, keyword := range b.Keywords {
		if keywordSet[keyword] {
			score += 1
		}
	}

	// Compare Category
	if a.Category == b.Category {
		score += 3
	}

	// Compare Author
	if a.Author == b.Author {
		score += 2
	}

	// Compare Site
	if a.Site == b.Site {
		score += 1
	}

	// Additional criteria can be added here

	return score
}

func addRSS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	db, err := initDB()
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

	// Slice to hold the URLs
	var urls []string

	// Iterate over the rows and append URLs to the slice
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			log.Fatal(err)
		}
		urls = append(urls, url)
	}

	// Check for any error encountered during iteration
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

	// Convert the response payload to JSON
	responseJSON, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, "Error creating response JSON", http.StatusInternalServerError)
		return
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON response
	w.Write(responseJSON)
}

func main() {
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	http.HandleFunc("/", streamAudio)
	http.HandleFunc("/embed", embedAudio)
	http.HandleFunc("/feed", feed)
	http.HandleFunc("/feed/add", addRSS)
	log.Println("Server started on port 443")
	log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/api.szn.io/fullchain.pem", "/etc/letsencrypt/live/api.szn.io/privkey.pem", nil))
}
