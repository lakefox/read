package main

import (
	"log"
	"net/http"

	"main/handlers"
)

func main() {

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./public"))))
	http.HandleFunc("/play", handlers.StreamAudio)
	http.HandleFunc("/embed", handlers.EmbedAudio)
	http.HandleFunc("/feed", handlers.Feed)
	http.HandleFunc("/feed/add", handlers.AddRSS)
	http.HandleFunc("/library/", handlers.HandleLibrary)

	log.Println("Server started on port 443")
	log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/api.szn.io/fullchain.pem", "/etc/letsencrypt/live/api.szn.io/privkey.pem", nil))
}
