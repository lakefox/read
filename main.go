package main

import (
	"log"
	"net/http"

	"main/handlers"
	"main/utils"
)

func main() {

	db, err := utils.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./public"))))
	http.HandleFunc("/play", handlers.StreamAudio(db))
	http.HandleFunc("/embed", handlers.EmbedAudio(db))
	http.HandleFunc("/feed", handlers.Feed(db))
	http.HandleFunc("/feed/add", handlers.AddRSS(db))
	http.HandleFunc("/library/", handlers.HandleLibrary(db))

	log.Println("Server started on port 443")
	log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/szn.io/fullchain.pem", "/etc/letsencrypt/live/szn.io/privkey.pem", nil))
}
