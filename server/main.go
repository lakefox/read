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
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	readability "github.com/go-shiori/go-readability"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

type Article struct {
	Keywords    []string
	Voice       string
	Category    string
	Title       string
	Author      string
	Image       string
	Description string
	Site        string
	Url         string
	Text        string
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
		Text:        fixtext.ReplaceDates(text),
	}, nil
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./audio.db")
	if err != nil {
		return nil, err
	}
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS audio_cache (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		category TEXT NOT NULL,
		keywords TEXT, -- Store keywords as a JSON string
		cf_url TEXT NOT NULL,
		title TEXT,
		author TEXT,
		image TEXT,
		description TEXT,
		site TEXT,
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

func uploadToR2(filePath, fileName string) (string, error) {
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	accessKeyID := os.Getenv("CLOUDFLARE_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("CLOUDFLARE_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("CLOUDFLARE_BUCKET_NAME")
	baseURL := os.Getenv("CLOUDFLARE_BASE_URL")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("auto"),
		Endpoint:    aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return "", err
	}

	uploader := s3.New(sess)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s", baseURL, fileName)
	return url, nil
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
	err = db.QueryRow("SELECT cf_url FROM audio_cache WHERE url = ?", queryURL).Scan(&cachedURL)
	if err == nil {
		http.Redirect(w, r, cachedURL, http.StatusFound)
		return
	}

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
	defer os.Remove(webmFile)

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

	uploadedURL, err := uploadToR2(webmFile, fmt.Sprintf("audio_%d.webm", time.Now().UnixNano()))
	if err != nil {
		log.Printf("Error uploading to Cloudflare R2: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	keywordsJSON, err := json.Marshal(opt.Keywords)
	if err != nil {
		log.Printf("Error marshaling keywords: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO audio_cache (url, category, keywords, cf_url, title, author, image, description, site) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		queryURL, opt.Category, string(keywordsJSON), uploadedURL, opt.Title, opt.Author, opt.Image, opt.Description, opt.Site)
	if err != nil {
		log.Printf("Error inserting into database: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
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
	err = db.QueryRow("SELECT url, category, keywords, cf_url, title, author, image, description, site FROM audio_cache WHERE url = ?", queryURL).Scan(
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
	}{
		AudioURL: audioURL,
		Image:    article.Image,
		Title:    article.Title,
		Author:   article.Author,
		Site:     article.Site,
		Category: article.Category,
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

func main() {
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	http.HandleFunc("/", streamAudio)
	http.HandleFunc("/embed", embedAudio)
	log.Println("Server started on port 443")
	log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/api.szn.io/fullchain.pem", "/etc/letsencrypt/live/api.szn.io/privkey.pem", nil))
}
