package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"

	// Replace with your actual module path
	"higher-or-lower/internal/game"
)

// PageData holds everything the index.tmpl needs
type PageData struct {
	Username   string
	Game       *game.GameState
	TopPlayers []game.Player
}

var (
	sessions  = make(map[string]*game.GameState)
	sessionMu sync.Mutex

	globalLeaderboard *game.Leaderboard
	indexTmpl         *template.Template
	loginTmpl         *template.Template
)

func init() {
	globalLeaderboard = game.LoadLeaderboard("leaderboard.json")
	indexTmpl = template.Must(template.ParseFiles("ui/html/index.tmpl"))
	loginTmpl = template.Must(template.ParseFiles("ui/html/login.tmpl"))
}

// getUsername checks the browser for our "player_name" cookie
func getUsername(r *http.Request) string {
	cookie, err := r.Cookie("player_name")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// LoginHandler displays the login form and sets the cookie
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

// LogoutHandler deletes the cookie
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "player_name",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// GameHandler is our main engine
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
				playerGame.Message = "please enter a valid number."
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
		TopPlayers: globalLeaderboard.GetTopPlayers(15),
	}
	indexTmpl.Execute(w, data)
}

// ResetHandler starts a new game for the logged-in user
func ResetHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsername(r)
	if username != "" {
		sessionMu.Lock()
		sessions[username] = game.New()
		sessionMu.Unlock()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
