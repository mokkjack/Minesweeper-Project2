/*
Prologue

Authors: Adam Berry, Barrett Brown, Jonathan Gott, Alex Phibbs, Minh Vu
Creation Date: 9/11/2025

Description:
- This file mostly handles the backend logic for a Minesweeper game.
It handles board structure, how each square would react to a certain event, bomb placement/generation,
flagging, recursive zero reveal/uncovering squares, win and lose conditions

Functions used:
- NewGameHandler: Creates a new game and board with bombs placed randomly on the board
	Input: number of mines
	Output: game handler with the board initialized

- AddNumbers: Adds the number to each square representing the number of adjacent bombs

- isiInbounds: Checks if a cell is inside the board

- GetBoard: Returns the state of the board

- RevealZero: Recursively uncovers zero-valued squares and their neighbors

- Click: Handles all clicks (user click, first click, lose/win, recursive uncovering)

- ToggleFlag: Toggles between flag states on a unrevealed square

- moveBombFrom: Changes the location of a bomb if the first click is a bomb

- revealAllBombs: Uncovers all bombs if it's in a lose condition

- checkWin: Check whether the game is in a win condition

Inputs:
- Board size
- Number of mines
- Clicks (clicks/flags)

Outputs:
- Updated game on clicks
- Game result (win/lose)
*/

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
	board      [][]Square
	rng        *rand.Rand
	firstClick bool
	gameOver   bool
	win        bool
	totalMines int
}

// This function should create the entire game board equipped with mines and numbered Squares
// Current iteration doesn't do that, but it does have a structure.
// Keep in mind we will have to do the generation of bombs, then seperately number the squares when implementing
func NewGameHandler(numMines int) Gamehandler {
	handler := Gamehandler{}
	handler.board = make([][]Square, config.BoardSize)
	handler.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	handler.firstClick = true
	handler.gameOver = false
	handler.win = false
	handler.totalMines = numMines

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

	// this slice will have all locations where mines can go
	possible_mine_locations := make([]int, 0, num_cells)

	// this for-loop finds every cell that is not the first clicked cell
	// and adds it to the list of possible mine locations
	for row := 0; row < config.BoardSize; row++ {
		for col := 0; col < config.BoardSize; col++ {
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

func (handler *Gamehandler) AddNumbers() {
	// For each square in the array, count the number of mines in the surrounding eight squares
	for row := 0; row < config.BoardSize; row++ {
		for col := 0; col < config.BoardSize; col++ {
			if handler.board[row][col].isBomb {
				handler.board[row][col].numValue = 0
				continue
			}
			bombc := 0
			for i := -1; i < 2; i++ {
				for j := -1; j < 2; j++ {
					if (i == 0 && j == 0) || !isiInbounds(handler, row+i, col+j) {
						continue
					}
					if handler.board[row+i][col+j].isBomb {
						bombc += 1
					}
				}
			}
			handler.board[row][col].numValue = bombc
		}
	}
}

func isiInbounds(handler *Gamehandler, row int, col int) bool {
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
	} else {
		sq.state = Uncovered
		return
	}

	// Recursively calls neighboring squares
	if sq.numValue == 0 {
		handler.RevealZero(row+1, col+1)
		handler.RevealZero(row+1, col-1)
		handler.RevealZero(row-1, col+1)
		handler.RevealZero(row-1, col-1)

		handler.RevealZero(row-1, col)
		handler.RevealZero(row+1, col)
		handler.RevealZero(row, col+1)
		handler.RevealZero(row, col-1)
	}
}

func (handler *Gamehandler) Click(row, col int) {
	if handler.gameOver || !isiInbounds(handler, row, col) {
		return
	}

	// First-click safety: if first click hits a bomb, move it elsewhere and recompute numbers.
	if handler.firstClick && handler.board[row][col].isBomb {
		handler.moveBombFrom(row, col)
	}
	handler.firstClick = false

	sq := &handler.board[row][col]
	if sq.state == Flagged || sq.state == Uncovered {
		return
	}

	if sq.isBomb {
		// lose
		handler.gameOver = true
		handler.win = false
		handler.revealAllBombs()
		return
	}

	if sq.numValue == 0 {
		handler.RevealZero(row, col)
	} else {
		sq.state = Uncovered
	}

	handler.checkWin()
}

// ToggleFlag flips flag state and checks win.
func (handler *Gamehandler) ToggleFlag(row, col int) {
	if handler.gameOver || !isiInbounds(handler, row, col) {
		return
	}
	sq := &handler.board[row][col]
	if sq.state == Uncovered {
		return
	}
	if sq.state == Flagged {
		sq.state = Covered
	} else {
		sq.state = Flagged
	}
	handler.checkWin()
}

// moveBombFrom relocates a bomb at (row,col) to the first safe non-bomb cell and re-runs AddNumbers.
func (handler *Gamehandler) moveBombFrom(row, col int) {
	handler.board[row][col].isBomb = false

	for r := 0; r < config.BoardSize; r++ {
		for c := 0; c < config.BoardSize; c++ {
			if (r == row && c == col) || handler.board[r][c].isBomb {
				continue
			}
			handler.board[r][c].isBomb = true
			handler.AddNumbers()
			return
		}
	}
	handler.board[row][col].isBomb = true
	handler.AddNumbers()
}

func (handler *Gamehandler) revealAllBombs() {
	for r := 0; r < config.BoardSize; r++ {
		for c := 0; c < config.BoardSize; c++ {
			if handler.board[r][c].isBomb {
				handler.board[r][c].state = Uncovered
			}
		}
	}
}

func (handler *Gamehandler) checkWin() {
	if handler.gameOver {
		return
	}
	flags := 0
	allNonBombsUncovered := true

	for r := 0; r < config.BoardSize; r++ {
		for c := 0; c < config.BoardSize; c++ {
			sq := handler.board[r][c]
			if sq.isBomb {
				if sq.state == Flagged {
					flags++
				}
			} else if sq.state != Uncovered {
				allNonBombsUncovered = false
			}
		}
	}
	// Win if every bomb is flagged and every non-bomb is uncovered
	if flags == handler.totalMines && allNonBombsUncovered {
		handler.gameOver = true
		handler.win = true
	}
}
