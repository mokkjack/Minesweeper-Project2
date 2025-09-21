
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
func NewGameHandler() Gamehandler {
	handler := Gamehandler{}
	handler.board = make([][]Square, config.BoardSize)

	for x := 0; x < config.BoardSize; x++ {
		handler.board[x] = make([]Square, config.BoardSize)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for row := 0; row < config.BoardSize; row++ {
		for col := 0; col < config.BoardSize; col++ {
			var box Square
			box.state = Covered
			box.isBomb = r.Intn(20) == 1 // currently no set number of bombs
			if !box.isBomb {
				box.numValue = r.Intn(2) // Placeholder for neighboring bomb count
			}
			handler.board[row][col] = box
		}
	}

	return handler
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
func (h *Gamehandler) mine_generator(first_click_row, first_click_col, num_mines int) {
	// find the cell that was the first click
	first_clicked_cell := first_click_row*config.BoardSize + first_click_col

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

			// If we come upon the first cell that was clicked, skip it so
			// that it does not become a possible mine location
			if cell_id != first_clicked_cell {
				possible_mine_locations = append(possible_mine_locations, cell_id)
			}
		}
	}

	// Shuffle the possible mine locations randomly
	h.rng.Shuffle(len(possible_mine_locations), func(i int, j int) {
		possible_mine_locations[i], possible_mine_locations[j] = possible_mine_locations[j], possible_mine_locations[i]
	})

	// Add mines to the first 10 or 20 cells in the shuffled list of possible mine locations
	for i := 0; i < num_mines; i++ {
		// get the cell location
		cell_id := possible_mine_locations[i]

		// convert cell id to row, column
		row := cell_id / config.BoardSize
		col := cell_id % config.BoardSize

		// add a mine to the cell
		h.board[row][col].isBomb = true
	}

}
