package internal

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Input represents the input data for the game
type Input struct {
	HorseLabel     string
	HorsesQuantity int
	ScoreTarget    int
	GameTimeout    string
}

// Horse represent the horse entity
type Horse struct {
	Label string
	Score atomic.Int32
}

var (
	horses      = []*Horse{}
	horsesMutex = &sync.RWMutex{}
)

// Winner greeting to the champion horse
func (h *Horse) Winner() string {
	return fmt.Sprintf("The horse winner is: %s - Score %d", h.Label, h.Score.Load())
}

// GetScore returns the current score safely
func (h *Horse) GetScore() int32 {
	return h.Score.Load()
}

// AddScore adds to the current score safely
func (h *Horse) AddScore(delta int32) {
	h.Score.Add(delta)
}
