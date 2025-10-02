//mediumAI.go
//Evan: 9/29/25 >> 1 hour
//Evan: 10/1/25 >> 2h
//Evan: 10/2/25 >> 3h

//func func_name(param_name param_type) return_type {}

//Components Package
package components

//Import Library
import (
	"fmt"
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
	// var picks []cell	//picks

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
		neighbor_tracker(handler, number_cells)
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

//Neighbor Tracker Function || nc = number_cells
func neighbor_tracker(handler *Gamehandler, nc []cell) []cell {
	//Local Variables
	var ntnc []cell
	next_to_number_cells := make([]cell, 0, config.BoardSize*config.BoardSize)
	
	//Call All Number Cells
	for i := 0; i < len(nc); i++ {
		ntnc = next_to_number_cells
		//Top-Left Cell
		if nc[i].r - 1 >= 0 && nc[i].r - 1 < config.BoardSize && nc[i].c - 1 >= 0 && nc[i].c - 1 < config.BoardSize {
			if handler.board[nc[i].r - 1][nc[i].c - 1].state == Uncovered && !(NTcontains(ntnc, cell{nc[i].r - 1, nc[i].c - 1})) {
				next_to_number_cells = append(next_to_number_cells, cell{nc[i].r - 1, nc[i].c - 1})
			}
		}
		//Top-Mid Cell
		if nc[i].r >= 0 && nc[i].r < config.BoardSize && nc[i].c - 1 >= 0 && nc[i].c - 1 < config.BoardSize {
			if handler.board[nc[i].r][nc[i].c - 1].state == Uncovered && !(NTcontains(ntnc, cell{nc[i].r, nc[i].c - 1})) {
				next_to_number_cells = append(next_to_number_cells, cell{nc[i].r, nc[i].c - 1})
			}
		}
		//Top-Right Cell
		if nc[i].r + 1 >= 0 && nc[i].r + 1 < config.BoardSize && nc[i].c - 1 >= 0 && nc[i].c - 1 < config.BoardSize {
			if handler.board[nc[i].r + 1][nc[i].c - 1].state == Uncovered && !(NTcontains(ntnc, cell{nc[i].r + 1, nc[i].c - 1})) {
				next_to_number_cells = append(next_to_number_cells, cell{nc[i].r + 1, nc[i].c - 1})
			}
		}
		//Mid-Left Cell
		if nc[i].r >= 0 && nc[i].r < config.BoardSize && nc[i].c - 1 >= 0 && nc[i].c - 1 < config.BoardSize {
			if handler.board[nc[i].r][nc[i].c - 1].state == Uncovered && !(NTcontains(ntnc, cell{nc[i].r, nc[i].c - 1})) {
				next_to_number_cells = append(next_to_number_cells, cell{nc[i].r, nc[i].c - 1})
			}
		}
		//Mid-Right Cell
		if nc[i].r >= 0 && nc[i].r < config.BoardSize && nc[i].c + 1 >= 0 && nc[i].c + 1 < config.BoardSize {
			if handler.board[nc[i].r][nc[i].c + 1].state == Uncovered && !(NTcontains(ntnc, cell{nc[i].r, nc[i].c + 1})) {
				next_to_number_cells = append(next_to_number_cells, cell{nc[i].r, nc[i].c + 1})
			}
		}
		//Bot-Left Cell
		if nc[i].r + 1 >= 0 && nc[i].r + 1 < config.BoardSize && nc[i].c - 1 >= 0 && nc[i].c - 1 < config.BoardSize {
			if handler.board[nc[i].r + 1][nc[i].c - 1].state == Uncovered && !(NTcontains(ntnc, cell{nc[i].r + 1, nc[i].c - 1})) {
				next_to_number_cells = append(next_to_number_cells, cell{nc[i].r + 1, nc[i].c - 1})
			}
		}
		//Bot-Mid Cell
		if nc[i].r + 1 >= 0 && nc[i].r + 1 < config.BoardSize && nc[i].c >= 0 && nc[i].c < config.BoardSize {
			if handler.board[nc[i].r + 1][nc[i].c].state == Uncovered && !(NTcontains(ntnc, cell{nc[i].r + 1, nc[i].c})) {
				next_to_number_cells = append(next_to_number_cells, cell{nc[i].r + 1, nc[i].c})
			}
		}
		//Bot-Right Cell
		if nc[i].r + 1 >= 0 && nc[i].r + 1 < config.BoardSize && nc[i].c + 1 >= 0 && nc[i].c + 1 < config.BoardSize {
			if handler.board[nc[i].r + 1][nc[i].c + 1].state == Uncovered && !(NTcontains(ntnc, cell{nc[i].r + 1, nc[i].c + 1})) {
				next_to_number_cells = append(next_to_number_cells, cell{nc[i].r + 1, nc[i].c + 1})
			}
		}
	}

	//Return Candidates
	fmt.Printf("%d\n", len(next_to_number_cells))
	return next_to_number_cells
}

//Neighbor Tracker Contains Function || tslice = this slice | tcell = this cell
func NTcontains(tslice []cell, tcell cell) bool {
	//Check all cells in current slice
	for i := 0; i < len(tslice); i++ {
		if tslice[i] == tcell { //Check for cell
			return true
		}
	}
	//No Identical Cell in Slice
	return false
}

/*
//
func checkSurrounding(row int, col int) cell {
	//Local Variable
	var priority int = 0


}

//
func countSurrounding(row int, col int) int {

}
*/




