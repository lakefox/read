package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
)

// RequestData represents the expected JSON structure for the POST request.
type RequestData struct {
	Text  string `json:"data"`
	Voice string `json:"voice"`
}

func streamAudio(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON request body
	var requestData RequestData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		log.Printf("Error decoding JSON: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var voice string
	if requestData.Voice == "female" {
		voice = "en_US-amy-medium.onnx"
	} else {
		voice = "en_US-joe-medium.onnx"
	}
	log.Printf("Using voice model: %s\n", voice)
	// Prepare the command using the text from the request
	piperCmd := exec.Command("sh", "-c", "echo '"+requestData.Text+"' | /root/server/piper/piper --model /root/"+voice+" --output-raw")
	ffmpegCmd := exec.Command("ffmpeg", "-f", "s16le", "-ar", "22050", "-ac", "1", "-i", "pipe:0", "-c:a", "libopus", "-f", "webm", "pipe:1")

	// Pipe the output of piper to ffmpeg
	piperOut, err := piperCmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe for piper: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ffmpegCmd.Stdin = piperOut

	// Set the HTTP headers for the response
	w.Header().Set("Content-Type", "audio/webm")

	// Pipe the output of ffmpeg to the HTTP response
	ffmpegOut, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe for ffmpeg: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Start the ffmpeg command first to ensure it's ready to receive data
	if err := ffmpegCmd.Start(); err != nil {
		log.Printf("Error starting ffmpeg command: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Start the piper command
	if err := piperCmd.Start(); err != nil {
		log.Printf("Error starting piper command: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Stream the ffmpeg output to the HTTP response
	go func() {
		if _, err := io.Copy(w, ffmpegOut); err != nil {
			log.Printf("Error streaming audio: %v\n", err)
		}
	}()

	// Wait for piper to finish
	if err := piperCmd.Wait(); err != nil {
		log.Printf("Error waiting for piper command: %v\n", err)
		// Ensure ffmpeg is terminated
		ffmpegCmd.Process.Kill()
		return
	}

	// Wait for ffmpeg to finish
	if err := ffmpegCmd.Wait(); err != nil {
		log.Printf("Error waiting for ffmpeg command: %v\n", err)
		return
	}
}

func main() {
	http.HandleFunc("/stream", streamAudio)

	log.Println("Server started on port 443")
	log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/tts.szn.io/fullchain.pem", "/etc/letsencrypt/live/tts.szn.io/privkey.pem", nil))
}
