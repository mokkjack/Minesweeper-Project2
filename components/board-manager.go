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

type gamehandler struct {
	board [][]Square
	rng   *rand.Rand
}

// This function should create the entire game board equipped with mines and numbered Squares
// Current iteration doesn't do that, but it does have a structure.
// Keep in mind we will have to do the generation of bombs, then seperately number the squares when implementing
func NewGameHandler() gamehandler {
	handler := gamehandler{}
	handler.board = make([][]Square, config.BoardSize)

	for x := 0; x < config.BoardSize; x++ {
		handler.board[x] = make([]Square, config.BoardSize)
	}
	rand.Seed(time.Now().UnixNano())
	for row := 0; row < config.BoardSize; row++ {
		for col := 0; col < config.BoardSize; col++ {
			var box Square
			box.state = Covered
			box.isBomb = rand.Intn(4) == 1 // currently no set number of bombs
			if !box.isBomb {
				box.numValue = rand.Intn(8)
			}
			handler.board[row][col] = box
		}
	}

	return handler
}

func GetBoard(handler *gamehandler) [][]Square {
	return handler.board
}

func (handler *gamehandler) RevealZero(row int, col int) {
	//Checks to see if cordinate is inside the board if not returns
	if row < 0 || row >= config.BoardSize || col < 0 || col >= config.BoardSize {
		return
	}

	sq := &handler.board[row][col] //Gets address at clicked position

	//If square is revealed or flagged don't reveal
	if sq.state == Uncovered || sq.state == Flagged {
		return
	}

	//If value is zero and not a bomb uncover
	if sq.numValue == 0 && !sq.isBomb {
		sq.state = Uncovered
	}

	//Recursively calls neighboring squares
	if sq.numValue == 0 {
		//Checks the rows
		for dr := -1; dr <= 1; dr++ {
			//Checks the columns
			for dc := -1; dc <= 1; dc++ {
				//Skips same spot
				if dr == 0 && dc == 0 {
					continue
				}
				handler.RevealZero(row+dr, col+dc)
			}
		}
	}
}
