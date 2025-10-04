// hardAI.go
// Ian Foehrweiser
// 10/3/2025
//

package components

import (
	"math/rand"
	"time"
	"minesweeper/config"
)

// Local cell struct for AI bookkeeping
type hardCell struct {
	r int
	c int
}

// HardAIMove SHOULD 1: Check safe moves, then the 1-2-1 rule, then randomly guess if it needs to
func HardAIMove(handler *Gamehandler) bool {
	//make sure game is running
	if handler == nil || handler.gameOver { 
		return false
	}

	// RNG initialization
	var rng *rand.Rand
	if handler.rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano())) // seed with current time
	} else {
		rng = handler.rng
	}

	// Collect covered and number cells
	coveredCells := make([]hardCell, 0, config.BoardSize*config.BoardSize)//hidden
	numberCells := make([]hardCell, 0, config.BoardSize*config.BoardSize) //uncovered

	for r := 0; r < config.BoardSize; r++ {
		for c := 0; c < config.BoardSize; c++ {
			sq := &handler.board[r][c]
			if sq.state == Covered {
				coveredCells = append(coveredCells, hardCell{r, c})
			} else if sq.state == Uncovered && sq.numValue > 0 {
				numberCells = append(numberCells, hardCell{r, c})
			}
		}
	}

	//no move
	if len(coveredCells) == 0 {
		return false
	}

	// count flagged neighbor cells
	for _, nc := range numberCells {
		neighbors := getCoveredNeighbors(handler, nc)
		if len(neighbors) == 0 {
			continue
		}

		// Count already placed flags
		flagCount := 0
		for _, neigh := range getAllNeighbors(handler, nc) {
			if handler.board[neigh.r][neigh.c].state == Flagged {
				flagCount++
			}
		}

		// Step 1: All remaining covered neighbors are safe if num == flagcount
		if handler.board[nc.r][nc.c].numValue == flagCount {
			move := neighbors[rng.Intn(len(neighbors))]
			handler.board[move.r][move.c].markedByAI = true
			handler.Click(move.r, move.c)
			return true
		}

		// Step 2: All covered neighbors are bombs if num == flagcount + hidden
		if handler.board[nc.r][nc.c].numValue == flagCount+len(neighbors) {
			move := neighbors[rng.Intn(len(neighbors))]
			handler.board[move.r][move.c].markedByAI = true
			handler.ToggleFlag(move.r, move.c)
			return true
		}
	}

	//  1-2-1 pattern rule 
	for r := 0; r < config.BoardSize; r++ {
		for c := 0; c < config.BoardSize-2; c++ {
			// Look for horizontally adjacent 1-2-1
			if handler.board[r][c].state == Uncovered &&
				handler.board[r][c+1].state == Uncovered &&
				handler.board[r][c+2].state == Uncovered &&
				handler.board[r][c].numValue == 1 &&
				handler.board[r][c+1].numValue == 2 &&
				handler.board[r][c+2].numValue == 1 {

				// Collect the cells above/below these three
				top := []hardCell{}
				bottom := []hardCell{}

				if r > 0 {
					// Above row
					if handler.board[r-1][c].state == Covered {
						top = append(top, hardCell{r - 1, c})
					}
					if handler.board[r-1][c+1].state == Covered {
						top = append(top, hardCell{r - 1, c + 1})
					}
					if handler.board[r-1][c+2].state == Covered {
						top = append(top, hardCell{r - 1, c + 2})
					}
				}

				if r < config.BoardSize-1 {
					// Below row
					if handler.board[r+1][c].state == Covered {
						bottom = append(bottom, hardCell{r + 1, c})
					}
					if handler.board[r+1][c+1].state == Covered {
						bottom = append(bottom, hardCell{r + 1, c + 1})
					}
					if handler.board[r+1][c+2].state == Covered {
						bottom = append(bottom, hardCell{r + 1, c + 2})
					}
				}

				// Apply 1-2-1 logic (check top first, then bottom)
				if len(top) == 3 {
					// Flag the two outer cells
					handler.ToggleFlag(top[0].r, top[0].c)
					handler.ToggleFlag(top[2].r, top[2].c)
					// Click the safe middle cell
					handler.Click(top[1].r, top[1].c)
					return true
				} else if len(bottom) == 3 {
					handler.ToggleFlag(bottom[0].r, bottom[0].c)
					handler.ToggleFlag(bottom[2].r, bottom[2].c)
					handler.Click(bottom[1].r, bottom[1].c)
					return true
				}
			}
		}
	}

	// --- Fallback: Guess randomly ---
	move := coveredCells[rng.Intn(len(coveredCells))]
	handler.board[move.r][move.c].markedByAI = true
	handler.Click(move.r, move.c)
	return true
}

// getCoveredNeighbors returns covered neighbors of a given number cell
func getCoveredNeighbors(handler *Gamehandler, nc hardCell) []hardCell {
	neighbors := make([]hardCell, 0, 8)
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			nr, nc2 := nc.r+dr, nc.c+dc
			if isInBounds(handler, nr, nc2) && handler.board[nr][nc2].state == Covered {
				neighbors = append(neighbors, hardCell{nr, nc2})
			}
		}
	}
	return neighbors
}

// getAllNeighbors returns all neighbors of a given number cell
func getAllNeighbors(handler *Gamehandler, nc hardCell) []hardCell {
	neighbors := make([]hardCell, 0, 8)
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			nr, nc2 := nc.r+dr, nc.c+dc
			if isInBounds(handler, nr, nc2) {
				neighbors = append(neighbors, hardCell{nr, nc2})
			}
		}
	}
	return neighbors
}

// isInBounds checks if row/col is in valid bounds
func isInBounds(handler *Gamehandler, r, c int) bool {
	return r >= 0 && r < config.BoardSize && c >= 0 && c < config.BoardSize
}