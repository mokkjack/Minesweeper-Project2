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

//Stack Structure - i dont want to hear it
type Stack struct {
	data	[]cell
	theSize	int
}

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
				covered_cucks = append(candidates, cell{r, c})
			}
		}
	}

	//Covered Cell Checker
	if len(candidates) == 0 {
		return false
	}

	//AI Move
	move := candidates[rng.Intn(len(candidates))]
	handler.Click(move.r, move.c)
	return true
	
}

//Narrow Possible Clicks Function
//func func_name(param_name param_type) return_type {}
func narrowSlice(candidates []cell) []cell {
	return candidates
}

