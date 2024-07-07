package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/server/plugins/category"
	fixtext "main/server/plugins/date"
	"main/server/plugins/gender"
	"main/server/plugins/keywords"
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
	_ "github.com/mattn/go-sqlite3"
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

func fetchDoc(url string) Article {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to parse %s, %v\n", url, err)
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
	}
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./audio.db")
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
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		image TEXT NOT NULL,
		description TEXT NOT NULL,
		site TEXT NOT NULL,
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
	err := godotenv.Load(".env")
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

	url := fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s", accountID, bucketName, fileName)
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

	if r.Method != http.MethodPost {
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

	opt := fetchDoc(queryURL)
	if err != nil {
		log.Printf("Error fetching text and voice from URL: %v\n", err)
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

	outputFile := fmt.Sprintf("/tmp/audio_%d.raw", time.Now().UnixNano())
	defer os.Remove(outputFile)

	piperCmd := exec.Command("sh", "-c", "echo '"+opt.Text+"' | /root/server/piper/piper --model /root/"+voiceModel+" --output-raw="+outputFile)
	ffmpegCmd := exec.Command("ffmpeg", "-f", "s16le", "-ar", "22050", "-ac", "1", "-i", outputFile, "-c:a", "libopus", "-f", "webm", "pipe:1")

	piperOut, err := piperCmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe for piper: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ffmpegCmd.Stdin = piperOut

	w.Header().Set("Content-Type", "audio/webm")

	ffmpegOut, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe for ffmpeg: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := ffmpegCmd.Start(); err != nil {
		log.Printf("Error starting ffmpeg command: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := piperCmd.Start(); err != nil {
		log.Printf("Error starting piper command: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	go func() {
		if _, err := io.Copy(w, ffmpegOut); err != nil {
			log.Printf("Error streaming audio: %v\n", err)
		}
	}()

	if err := piperCmd.Wait(); err != nil {
		log.Printf("Error waiting for piper command: %v\n", err)
		ffmpegCmd.Process.Kill()
		return
	}

	if err := ffmpegCmd.Wait(); err != nil {
		log.Printf("Error waiting for ffmpeg command: %v\n", err)
		return
	}

	uploadedURL, err := uploadToR2(outputFile, fmt.Sprintf("audio_%d.webm", time.Now().UnixNano()))
	if err != nil {
		log.Printf("Error uploading to Cloudflare R2: %v\n", err)
		return
	}
	keywordsJSON, err := json.Marshal(opt.Keywords)
	if err != nil {
		fmt.Println("Error marshaling keywords:", err)
		return
	}
	_, err = db.Exec("INSERT INTO audio_cache (url, category, keywords, cf_url) VALUES (?, ?, ?, ?)", queryURL, opt.Category, string(keywordsJSON), uploadedURL)
	if err != nil {
		log.Printf("Error inserting into database: %v\n", err)
	}
}

func main() {
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	http.HandleFunc("/", streamAudio)
	// http.HandleFunc("/player", streamAudio)
	log.Println("Server started on port 443")
	log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/tts.szn.io/fullchain.pem", "/etc/letsencrypt/live/tts.szn.io/privkey.pem", nil))
}
