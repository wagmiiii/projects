package game

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
)

// --- GAME LOGIC ---

type GameState struct {
	Difficulty     string
	MaxNumber      int
	Target         int
	GuessesAllowed int
	GuessesTaken   int
	Guesses        []int
	Message        string
	GameOver       bool
}

func New(difficulty string) *GameState {
	maxNumber := 100
	guessesAllowed := 5

	switch difficulty {
	case "Easy":
		maxNumber = 50
		guessesAllowed = 5
	case "Hard":
		maxNumber = 500
		guessesAllowed = 8
	default:
		difficulty = "Standard"
		maxNumber = 100
		guessesAllowed = 5
	}

	return &GameState{
		Difficulty:     difficulty,
		MaxNumber:      maxNumber,
		Target:         rand.Intn(maxNumber + 1),
		GuessesAllowed: guessesAllowed,
		GuessesTaken:   0,
		Guesses:        []int{},
		Message:        "Welcome! Make your first guess.",
		GameOver:       false,
	}
}

func (g *GameState) CheckGuess(guess int) {
	g.GuessesTaken++
	g.Guesses = append(g.Guesses, guess)
	guessesLeft := g.GuessesAllowed - g.GuessesTaken

	if guess == g.Target {
		g.Message = "BIM!! Good job, you guessed it!"
		g.GameOver = true
	} else if g.GuessesTaken >= g.GuessesAllowed {
		g.Message = fmt.Sprintf("Chaii, so sorry, the answer was %d.", g.Target)
		g.GameOver = true
	} else if guess < g.Target {
		g.Message = fmt.Sprintf("Oh snap!! The number is higher! You have %d guesses left.", guessesLeft)
	} else {
		g.Message = fmt.Sprintf("Oh snap!! The number is lower! You have %d guesses left.", guessesLeft)
	}
}

// --- LEADERBOARD LOGIC ---

type Player struct {
	Username    string `json:"username"`
	GamesPlayed int    `json:"gamesPlayed"`
	Wins        int    `json:"wins"`
	Losses      int    `json:"losses"`
	WinRatio    int    `json:"winRatio"` // Catch the generated percentage
}

type Leaderboard struct {
	DB *sql.DB
}

// RecordGame saves the result to Supabase
func (lb *Leaderboard) RecordGame(username string, won bool) {
	winsToAdd := 0
	lossesToAdd := 0
	if won {
		winsToAdd = 1
	} else {
		lossesToAdd = 1
	}

	query := `
		INSERT INTO leaderboard (username, games_played, wins, losses)
		VALUES ($1, 1, $2, $3)
		ON CONFLICT (username)
		DO UPDATE SET
			games_played = leaderboard.games_played + 1,
			wins = leaderboard.wins + $2,
			losses = leaderboard.losses + $3;
	`

	_, err := lb.DB.Exec(query, username, winsToAdd, lossesToAdd)
	if err != nil {
		log.Println("Error saving score to database:", err)
	}
}

// GetTopPlayers pulls the leaderboard, including the generated win_ratio
func (lb *Leaderboard) GetTopPlayers(limit int) []Player {
	query := `SELECT username, games_played, wins, losses, win_ratio FROM leaderboard ORDER BY win_ratio DESC LIMIT $1`

	rows, err := lb.DB.Query(query, limit)
	if err != nil {
		log.Println("Error fetching leaderboard:", err)
		return []Player{}
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.Username, &p.GamesPlayed, &p.Wins, &p.Losses, &p.WinRatio); err == nil {
			players = append(players, p)
		}
	}
	return players
}
