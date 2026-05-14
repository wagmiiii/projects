package game

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"
)

// GameState holds the state of a single game.
type GameState struct {
	Target         int
	GuessesAllowed int
	GuessesTaken   int
	Guesses        []int
	Message        string
	GameOver       bool
}

// New initializes and returns a pointer to a fresh GameState.
func New() *GameState {
	return &GameState{
		Target:         rand.Intn(101),
		GuessesAllowed: 5,
		GuessesTaken:   0,
		Guesses:        []int{},
		Message:        "welcome! make your first guess.",
		GameOver:       false,
	}
}

// CheckGuess evaluates the user's input and updates the game state.
func (g *GameState) CheckGuess(guess int) {
	g.GuessesTaken++
	g.Guesses = append(g.Guesses, guess)
	guessesLeft := g.GuessesAllowed - g.GuessesTaken

	if guess == g.Target {
		g.Message = "good job, you guessed it!"
		g.GameOver = true
	} else if g.GuessesTaken >= g.GuessesAllowed {
		g.Message = fmt.Sprintf("chaii, so sorrry. the answer was %d.", g.Target)
		g.GameOver = true
	} else if guess < g.Target {
		g.Message = fmt.Sprintf("the number is higher! You have %d guesses left.", guessesLeft)
	} else {
		g.Message = fmt.Sprintf("the number is lower! You have %d guesses left.", guessesLeft)
	}
}

// --- LEADERBOARD LOGIC ---

// Player holds the stats for a single user
type Player struct {
	Username    string  `json:"username"`
	GamesPlayed int     `json:"gamesPlayed"`
	Wins        int     `json:"wins"`
	Losses      int     `json:"losses"`
	WinRatio    float64 `json:"winRatio"`
}

// Leaderboard manages the map of players and file saving.
type Leaderboard struct {
	Players map[string]Player `json:"players"`
	mu      sync.Mutex
}

// LoadLeaderboard reads the JSON file and turns it into Go structs
func LoadLeaderboard(filename string) *Leaderboard {
	lb := &Leaderboard{
		Players: make(map[string]Player),
	}

	fileData, err := os.ReadFile(filename)
	if err != nil {
		return lb // File doesn't exist yet, return empty leaderboard
	}

	json.Unmarshal(fileData, &lb.Players)
	return lb
}

// Save writes the current Go map back into the JSON file
func (lb *Leaderboard) Save(filename string) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	jsonData, err := json.MarshalIndent(lb.Players, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

// RecordGame updates a player's stats and saves to the file
func (lb *Leaderboard) RecordGame(username string, won bool) {
	lb.mu.Lock()

	player, exists := lb.Players[username]
	if !exists {
		player = Player{Username: username}
	}

	player.GamesPlayed++
	if won {
		player.Wins++
		player.WinRatio = math.Round(float64(player.Wins) / float64(player.GamesPlayed) * 100)
	} else {
		player.Losses++
		player.WinRatio = math.Round(float64(player.Wins) / float64(player.GamesPlayed) * 100)
	}

	lb.Players[username] = player
	lb.mu.Unlock()

	lb.Save("leaderboard.json")
}

// GetTopPlayers returns a sorted slice of the best players
func (lb *Leaderboard) GetTopPlayers(limit int) []Player {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	var players []Player
	for _, p := range lb.Players {
		players = append(players, p)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].WinRatio > players[j].WinRatio
	})

	if len(players) > limit {
		return players[:limit]
	}
	return players
}
