package game

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
)

// --- GAME LOGIC ---

type GameState struct {
	Difficulty   string
	MaxNumber    int
	Target       int
	GuessesTaken int
	Guesses      []int
	Message      string
	GameOver     bool
}

func New(difficulty string) *GameState {
	maxNumber := 100

	switch difficulty {
	case "Easy":
		maxNumber = 50
	case "Hard":
		maxNumber = 500
	default:
		difficulty = "Standard"
		maxNumber = 100
	}

	return &GameState{
		Difficulty:   difficulty,
		MaxNumber:    maxNumber,
		Target:       rand.Intn(maxNumber + 1),
		GuessesTaken: 0,
		Guesses:      []int{},
		Message:      "Welcome! Make your first guess.",
		GameOver:     false,
	}
}

func (g *GameState) CheckGuess(guess int) {
	g.GuessesTaken++
	g.Guesses = append(g.Guesses, guess)

	if guess == g.Target {
		g.Message = fmt.Sprintf("BIM!! Good job, you guessed it in %d tries!", g.GuessesTaken)
		g.GameOver = true
	} else if guess < g.Target {
		g.Message = "Oh snap!! The number is higher!"
	} else {
		g.Message = "Oh snap!! The number is lower!"
	}
}

// --- LEADERBOARD LOGIC ---

type Player struct {
	Username    string `json:"username"`
	GamesPlayed int    `json:"gamesPlayed"`
	BestScore   int    `json:"bestScore"` // Fewest tries taken
}

type Leaderboard struct {
	DB *sql.DB
}

// RecordWin updates the player's total games and their best (lowest) score.
func (lb *Leaderboard) RecordWin(username string, triesTaken int) {
	query := `
		INSERT INTO leaderboard (username, games_played, best_score)
		VALUES ($1, 1, $2)
		ON CONFLICT (username)
		DO UPDATE SET
			games_played = leaderboard.games_played + 1,
			best_score = LEAST(leaderboard.best_score, $2);
	`

	_, err := lb.DB.Exec(query, username, triesTaken)
	if err != nil {
		log.Println("Error saving score to database:", err)
	}
}

// GetTopPlayers pulls the leaderboard ordered by the lowest best_score.
func (lb *Leaderboard) GetTopPlayers(limit int) []Player {
	// ASC order because fewer tries is better
	query := `SELECT username, games_played, best_score FROM leaderboard ORDER BY best_score ASC LIMIT $1`

	rows, err := lb.DB.Query(query, limit)
	if err != nil {
		log.Println("Error fetching leaderboard:", err)
		return []Player{}
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.Username, &p.GamesPlayed, &p.BestScore); err == nil {
			players = append(players, p)
		}
	}
	return players
}
