package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func DownloadAudio(db *sql.DB, url string) error {
	opt, err := FetchDoc(url)

	fmt.Println(opt)
	if err != nil {
		return fmt.Errorf("error fetching document: %v", err)
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

	piperCmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' | /root/piper/piper --model /root/%s --output-raw", opt.Text, voiceModel))
	ffmpegCmd := exec.Command("ffmpeg", "-f", "s16le", "-ar", "22050", "-ac", "1", "-i", "pipe:0", "-c:a", "libopus", "-f", "webm", webmFile)

	piperStdout, err := piperCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error getting stdout pipe: %v", err)
	}

	ffmpegCmd.Stdin = piperStdout

	if err := piperCmd.Start(); err != nil {
		return fmt.Errorf("error starting Piper command: %v", err)
	}

	if err := ffmpegCmd.Start(); err != nil {
		return fmt.Errorf("error starting FFmpeg command: %v", err)
	}

	if err := piperCmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for Piper command: %v", err)
	}

	if err := ffmpegCmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for FFmpeg command: %v", err)
	}

	keywordsJSON, err := json.Marshal(opt.Keywords)
	if err != nil {
		return fmt.Errorf("error marshaling keywords: %v", err)
	}

	_, err = db.Exec("INSERT INTO audio_cache (url, category, keywords, location, title, author, image, description, site, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, date)",
		url, opt.Category, string(keywordsJSON), webmFile, opt.Title, opt.Author, opt.Image, opt.Description, opt.Site, opt.Date)
	if err != nil {
		return fmt.Errorf("error inserting into database: %v", err)
	}

	return nil
}

func StreamNewAudio(w http.ResponseWriter, r *http.Request, db *sql.DB, queryURL string) {
	opt, err := FetchDoc(queryURL)
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
