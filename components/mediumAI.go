//mediumAI.go
//Evan: 9/29/25 >> 1 hour
//Evan: 10/1/25 >> 2h
//Evan: 10/2/25 >> 1h

//func func_name(param_name param_type) return_type {}

//Components Package
package components

//Import Library
import (
	"math/rand"
	"time"
	"minesweeper/config"
)

//Stack Data Structure - dont even say anything
type Stack struct {
	data		[]cell
	theSize		int
	theCapacity	int
}

//Cell Structure
type cell struct {
	r	int
	c	int
}

//Medium AI Move Function
func MediumAIMove(handler *Gamehandler) bool {
	//Game Condition Checker
	if handler == nil || handler.gameOver {
		return false
	}

	//Local Variables
	var rng *rand.Rand 	//random
	var guess bool		//toggle guess

	//Initialize rng if it not created
	if handler.rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else { 
		rng = handler.rng
	}

	//Initialize Stack
	// theStack := Stack{}

	//Collect All Covered Cells
	covered_cells := make([]cell, 0, config.BoardSize*config.BoardSize)
	number_cells := make([]cell, 0, config.BoardSize*config.BoardSize)

	for r := 0; r < config.BoardSize; r++ {
		for c := 0; c < config.BoardSize; c++ {
			sq := &handler.board[r][c]
			if sq.state == Covered {
				covered_cells = append(covered_cells, cell{r, c})
			} else if sq.state == Uncovered && sq.numValue != 0 {
				number_cells = append(number_cells, cell{r, c})
			}
		}
	}

	//
	if len(number_cells) == 0 {
		guess = true
	} else {
		// candidates = splitCover(number_cells)
		guess = false
	}

	//Covered Cell Checker
	if len(covered_cells) == 0 {
		return false
	}

	//AI Move Decider
	if guess {
		move := covered_cells[rng.Intn(len(covered_cells))]
		print("guess")

		//Highlight and Make AI Move
		handler.board[move.r][move.c].markedByAI = true
		handler.Click(move.r, move.c)
		return true

	} else {
		move := covered_cells[rng.Intn(len(covered_cells))] //needs to be adjusted
		print("not guess")

		//Highlight and Make AI Move
		handler.board[move.r][move.c].markedByAI = true
		handler.Click(move.r, move.c)
		return true
	}
}

//
func splitCover(number_cells []cell) []cell {
	// next_to_number_cells := make([]cell, 0, config.BoardSize*config.BoardSize)

	return number_cells
}




