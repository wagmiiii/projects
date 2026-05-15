package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"

	// Replace with your actual module path
	"github.com/wagmiiii/projects/internal/game"
)

type PageData struct {
	Username   string
	Game       *game.GameState
	TopPlayers []game.Player
}

var (
	sessions   = make(map[string]*game.GameState)
	sessionMu  sync.Mutex
	
	globalLeaderboard *game.Leaderboard
	indexTmpl         *template.Template
	loginTmpl         *template.Template
)

// InitHandlers connects the handlers to the DB and loads templates
func InitHandlers(db *sql.DB) {
	globalLeaderboard = &game.Leaderboard{DB: db}
	indexTmpl = template.Must(template.ParseFiles("ui/html/index.tmpl"))
	loginTmpl = template.Must(template.ParseFiles("ui/html/login.tmpl"))
}

func getUsername(r *http.Request) string {
	cookie, err := r.Cookie("player_name")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		loginTmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		if username != "" {
			http.SetCookie(w, &http.Cookie{
				Name:    "player_name",
				Value:   username,
				Expires: time.Now().Add(24 * time.Hour),
				Path:    "/",
			})
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "player_name",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsername(r)
	if username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sessionMu.Lock()
	playerGame, exists := sessions[username]
	if !exists {
		playerGame = game.New()
		sessions[username] = playerGame
	}
	sessionMu.Unlock()

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err == nil {
			guessStr := r.FormValue("guess")
			guess, err := strconv.Atoi(guessStr)
			if err != nil {
				playerGame.Message = "Please enter a valid number."
			} else {
				playerGame.CheckGuess(guess)

				if playerGame.GameOver {
					won := (guess == playerGame.Target)
					globalLeaderboard.RecordGame(username, won)
				}
			}
		}
	}

	data := PageData{
		Username:   username,
		Game:       playerGame,
		TopPlayers: globalLeaderboard.GetTopPlayers(50), // Top 5
	}
	indexTmpl.Execute(w, data)
}

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	username := getUsername(r)
	if username != "" {
		sessionMu.Lock()
		sessions[username] = game.New()
		sessionMu.Unlock()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}