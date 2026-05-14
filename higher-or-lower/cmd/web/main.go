package main

import (
	"log"
	"net/http"

	"higher-or-lower/internal/handlers"
)

func main() {
	// Serve static files (CSS, images, JS) from the ui/static directory
	// 1. Point to the folder where the static files live
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// 2. Strip the "/static" prefix from the URL, and serve the files
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Map the URL paths to the exported functions in the handlers package
	http.HandleFunc("/", handlers.GameHandler)
	http.HandleFunc("/reset", handlers.ResetHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	log.Println("Server starting on http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}