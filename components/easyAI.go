// Zhang: easy AI mode
package components

import (
	"math/rand"
	"time"

	"minesweeper/config"
)

// function for easy AI
func EasyAIMove(handler *Gamehandler) bool {
	//check game condition
	if handler == nil || handler.gameOver {
		return false
	}
	//declare variable for random
	var rng *rand.Rand
	//created rng if it not created
	if handler.rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
		//if there is one, use that created one instead
	} else {
		rng = handler.rng
	}
	//defined a structure to store row and column coordinate
	type cell struct{ r, c int }

	// Collect all covered cells (potential moves) into a slice
	candidates := make([]cell, 0, config.BoardSize*config.BoardSize)
	for r := 0; r < config.BoardSize; r++ {
		for c := 0; c < config.BoardSize; c++ {
			sq := &handler.board[r][c]
			if sq.state == Covered {
				candidates = append(candidates, cell{r, c})
			}
		}
	}
	//if there is no covered cell return false
	if len(candidates) == 0 {
		return false
	}

	// Randomly select a candidate cell
	move := candidates[rng.Intn(len(candidates))]
	//perform the click function
	handler.Click(move.r, move.c)
	//move is successfully made
	return true
}
