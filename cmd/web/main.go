package main

import (
	"log"
	"net/http"
	"time"
	"os"
	"database/sql"

	// The PostgreSQL driver
	_ "github.com/lib/pq"

	"github.com/wagmiiii/projects/internal/handlers"
)

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now
		log.Printf("%s %s | %s", r.Method, r.URL, time.Since(start()))
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 1. Connect to Supabase
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	log.Println("Successfully connected to Supabase!")

	// 2. Initialize Handlers with the DB
	handlers.InitHandlers(db)

	// Serve static files (CSS, images, JS) from the ui/static directory
	// 1. Point to the folder where the static files live
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// 2. Strip the "/static" prefix from the URL, and serve the files
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Map the URL paths to the exported functions in the handlers package
	mux := http.NewServeMux()

	mux.HandleFunc("/", loggingMiddleware(handlers.GameHandler))
	mux.HandleFunc("/reset", loggingMiddleware(handlers.ResetHandler))
	mux.HandleFunc("/login", loggingMiddleware(handlers.LoginHandler))
	mux.HandleFunc("/logout", loggingMiddleware(handlers.LogoutHandler))

	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local development
	}

	log.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
