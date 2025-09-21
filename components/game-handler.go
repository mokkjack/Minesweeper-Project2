
package components

import (
	"math/rand"
	"minesweeper/config"
	"time"
)

type SquareState int

const (
	Covered SquareState = iota
	Uncovered
	Flagged
)

// Define the square struct
type Square struct {
	state    SquareState
	isBomb   bool
	numValue int
}

type Gamehandler struct {
	board [][]Square
	rng   *rand.Rand
}

// This function should create the entire game board equipped with mines and numbered Squares
// Current iteration doesn't do that, but it does have a structure.
// Keep in mind we will have to do the generation of bombs, then seperately number the squares when implementing
func NewGameHandler(numMines int) Gamehandler {
	handler := Gamehandler{}
	handler.board = make([][]Square, config.BoardSize)
	handler.rng = rand.New(rand.NewSource(time.Now().UnixNano()))


	for x := 0; x < config.BoardSize; x++ {
		handler.board[x] = make([]Square, config.BoardSize)
	}

	// Iterate through and initialize a square struct for each index in the array
	for row := 0; row < config.BoardSize; row++ {
		for col := 0; col < config.BoardSize; col++ {
			var box Square
			box.state = Covered
			handler.board[row][col] = box
		}
	}
	

	// represents the total number of cells
	num_cells := config.BoardSize * config.BoardSize

	// this slice will have all locations where mines CAN go
	possible_mine_locations := make([]int, 0, num_cells)

	// this for-loop finds every cell that is not the first clicked cell
	// and adds it to the list of possible mine locations
	for row := 0; row < config.BoardSize; row++ {
		for col := 0; col <= config.BoardSize; col++ {
			// find current cell
			cell_id := row*config.BoardSize + col

			possible_mine_locations = append(possible_mine_locations, cell_id)
		}
	}

	// Shuffle the possible mine locations randomly
	handler.rng.Shuffle(len(possible_mine_locations), func(i int, j int) {
		possible_mine_locations[i], possible_mine_locations[j] = possible_mine_locations[j], possible_mine_locations[i]
	})

	// Add mines to the first 10 or 20 cells in the shuffled list of possible mine locations
	for i := 0; i < numMines; i++ {
		// get the cell location
		cell_id := possible_mine_locations[i]

		// convert cell id to row, column
		row := cell_id / config.BoardSize
		col := cell_id % config.BoardSize

		// add a mine to the cell
		handler.board[row][col].isBomb = true
	}

	handler.AddNumbers()

	return handler
}

func (handler *Gamehandler) AddNumbers(){
	// For each square in the array, count the number of mines in the surrounding eight squares
	for row := 0; row < config.BoardSize; row++ { 
		for col := 0; col < config.BoardSize; col++{
			if handler.board[row][col].isBomb{
				continue
			}
			bombc := 0
			for i := -1; i < 2; i++{
				for j := -1; j < 2; j++{
					if (i == 0 && j == 0) || !isiInbounds(handler, row + i, col+j){
						continue
					}
					if(handler.board[row+i][col+j].isBomb){
						bombc += 1
					}
				}
			}
			handler.board[row][col].numValue = bombc
		}
	}
	for x := 0; x < config.BoardSize; x++ {
		for y := 0; y < config.BoardSize; y++{
			if(handler.board[x][y].isBomb){
				print("B ")
			} else {
				print(handler.board[x][y].numValue)
				print(" ")
			}
		}
		println()
	}
}

func isiInbounds(handler *Gamehandler, row int, col int) bool{
	return (row >= 0) && (row < config.BoardSize) && (col >= 0) && (col < config.BoardSize)
}

func GetBoard(handler *Gamehandler) [][]Square {
	return handler.board
}

func (handler *Gamehandler) RevealZero(row int, col int) {
	// Checks to see if coordinate is inside the board if not returns
	if row < 0 || row >= config.BoardSize || col < 0 || col >= config.BoardSize {
		return
	}

	sq := &handler.board[row][col] //Gets address at clicked position

	// If square is revealed or flagged don't reveal
	if sq.state == Uncovered || sq.state == Flagged {
		return
	}

	// If value is zero and not a bomb uncover
	if sq.numValue == 0 && !sq.isBomb {
		sq.state = Uncovered
	}

	// Recursively calls neighboring squares
	if sq.numValue == 0 {
		handler.RevealZero(row+1, col)
		handler.RevealZero(row-1, col)
		handler.RevealZero(row, col+1)
		handler.RevealZero(row, col-1)
	}
}


// mine generator functions needs the location of the first click, and the number of mines
// func (h *Gamehandler) mine_generator(first_click_row, first_click_col, num_mines int) {
