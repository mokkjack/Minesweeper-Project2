//hardAI.go
//Evan: 9/29/25 >> 1 hour

// Import Libaries
package components

//Import Library
import (
	"math/rand"
	"minesweeper/config"
	"time"
)

// Function
func HardAIMove(handler *Gamehandler) bool {
	if handler == nil || handler.gameOver {
		return false
	}

	var rng *rand.Rand

	if handler.rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		rng = handler.rng
	}

	type cell struct{ r, c int }

	candidates := make([]cell, 0, config.BoardSize*config.BoardSize)
	for r := 0; r < config.BoardSize; r++ {
		for c := 0; c < config.BoardSize; c++ {
			sq := &handler.board[r][c]
			if sq.state == Covered {
				candidates = append(candidates, cell{r, c})
			}
		}
	}

	if len(candidates) == 0 {
		return false
	}

	move := candidates[rng.Intn(len(candidates))]
	handler.Click(move.r, move.c)
	return true

}
