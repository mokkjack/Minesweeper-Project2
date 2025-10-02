//mediumAI.go
//Evan: 9/29/25 >> 1 hour
//Evan: 10/1/25 >> 2h

//Components Package
package components

//Import Library
import (
	"math/rand"
	"time"
	"minesweeper/config"
)

//Medium AI Move Function
func MediumAIMove(handler *Gamehandler) bool {
	//Game Condition Checker
	if handler == nil || handler.gameOver {
		return false
	}

	//Random Variable
	var rng *rand.Rand

	//Initialize rng if it not created
	if handler.rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else { 
		rng = handler.rng
	}

	//Coordinate Storage Struct Initalization
	type cell struct { r, c int }

	//Collect All Covered Cells
	covered_cucks := make([]cell, 0, config.BoardSize*config.BoardSize)
	numbered_cucks := make([]cell, 0, config.BoardSize*config.BoardSize)

	for r := 0; r < config.BoardSize; r++ {
		for c := 0; c < config.BoardSize; c++ {
			sq := &handler.board[r][c]
			if sq.state == Covered {
				covered_cucks = append(covered_cucks, cell{r, c})
			} else if sq.state == Uncovered {
				numbered_cucks = append(numbered_cucks, cell{r, c})
			}
		}
	}

	//Covered Cell Checker
	if len(covered_cucks) == 0 {
		return false
	}

	//AI Move
	move := covered_cucks[rng.Intn(len(covered_cucks))] //NEEDS TO BE ADJUSTED

	//Highlight AI Move
	handler.board[move.r][move.c].markedByAI = true
	handler.Click(move.r, move.c)
	return true
	
}

//Narrow Possible Clicks Function
//func func_name(param_name param_type) return_type {}
func narrowSlice(blankCells []cell, numCells []cell) []cell {
	return blankCells
}

//Check Surrounding 
func checkSurround(curr_cell cell) int {
	return 0
}

