/*
Prologue

Authors: Adam Berry, Barrett Brown, Jonathan Gott, Alex Phibbs, Minh Vu
Creation Date: 9/11/2025

Description:
- This file mostly handles the backend logic for a Minesweeper game.
It handles board structure, how each square would react to a certain event, bomb placement/generation,
flagging, recursive zero reveal/uncovering squares, win and lose conditions

Functions:
- NewGameHandler: Creates a new game and board with bombs placed randomly on the board
	Input: number of mines
	Output: game handler with the board initialized

- AddNumbers: Makes the number of each square equal to the number representing the adjacent bombs

- isiInbounds: Helper function, checks if a cell is inside the board

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
  - Does this through updating the underlyining "Square" struct so other functions given this context can properly tell what's been changed
- Game result (win/lose)
*/

package components

import (
	"fmt"
	"math/rand"
	"minesweeper/config"
	"time"
)

type SquareState int

// Constant used to represent state so we can call to this to figure out if something is covered/flagged or waht not to properly reflect elsewhere
const (
	Covered SquareState = iota
	Uncovered
	Flagged
)

// Define the square struct, this is used for the cells in ui-handler.go but allows you to see cell state/if cell=bomb and the number of neighbors that cell has (if not bomb)
type Square struct {
	state      SquareState // If something is covered/uncovered/flagged
	isBomb     bool        // If something is a bomb
	numValue   int         // Neighbor count
	markedByAI bool        // Whether the square was clicked by the AI
}

// Gamehandler structs holds the board sets the rng value and whether this is firstclick and if the game is over (win or not) and the total number of mines
type Gamehandler struct {
	board      [][]Square // Used to store underlyining board
	rng        *rand.Rand // Used for bomb generation
	firstClick bool       // Used to ensure if this is first click + bomb we dont insta lose
	gameOver   bool       // Used to ensure no more game/also to trigger win/lost message
	win        bool       // Used to tell ui-handler to show win/lost
	totalMines int        // Used in NewGameHandler

	//Zhang: turn-based AI support
	aiEnabled    bool   // Whether AI is enabled
	aiTurn       bool   // Whether it's AI's turn
	aiDifficulty string // use for diffculty selection
	aiSolver     bool   // Whether Solver mode is enabled
}

// This function creates the game board equipped with mines and numbered squares
// Inputs: numMines as an int to place on the board
// Outputs: A gamehandler struct so you can adjust/look at the board
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

	// Called to adjust the "neighbor numbers" of each cell
	handler.AddNumbers()

	return handler
}

// Function that iterates through the game board and counts all nearby cells and sees how many bombs there are and sets it's numValue equal to that
// Inputs: handler object containing the game board
// Outputs: None, adjusts the underlining handler object
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

// Helper function for finding if a row/col is in bounds based on a game bound
// Inputs: Row/Col and handler object for game board
// Outputs: Bool value representing if in bounds
func isiInbounds(handler *Gamehandler, row int, col int) bool {
	return (row >= 0) && (row < config.BoardSize) && (col >= 0) && (col < config.BoardSize)
}

// Helper function to get the board of the handler object specifically
// Inputs: Handler object
// Outputs: [][]Square Game board object
func GetBoard(handler *Gamehandler) [][]Square {
	return handler.board
}

// Recursive function that "floods" the spot clicked revealing all up to the first "number" value
// Inputs: Handler, row and col value
// Outputs: None, updates state on the square
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

// Function that handles everything the click needs to do from first click safety to bomb discovered to calling recursive flood function/win condition
// Inputs: Row/Col and game handler object
// Outputs: None, ensures proper representation on the 2D-array as well as ending the game if need be by calling win codition
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
	UpdateGameUI(handler)
	handler.checkWin()
}

// ToggleFlag flips flag state and checks win.
// Inputs: row/col and gamehandler object
// Outputs: Nothing just edits the flagged state
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
		if handler.aiEnabled {
			handler.aiTurn = true
		}
	} else {
		sq.state = Flagged
		if handler.aiEnabled {
			handler.aiTurn = true
		}
	}
	handler.checkWin()
}

// Function that relocates a bomb at (row,col) to the first safe non-bomb cell and re-runs AddNumbers.
// Inputs: gameHandler object and row/col
// Outputs: Nothing just regenerates board into a safe "first-click" state
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

// Function that upon losing will be called, just iterates through the cells and if it is a bomb reveals it
// Inputs: gameHandler object
// Outputs: None, just edits the board
func (handler *Gamehandler) revealAllBombs() {
	for r := 0; r < config.BoardSize; r++ {
		for c := 0; c < config.BoardSize; c++ {
			if handler.board[r][c].isBomb {
				handler.board[r][c].state = Uncovered
			}
		}
	}
}

// Funciton used to check win condition (if all non-bombs uncovered)
// TODO: Could edit out flag variable and stuff just had that as I thought win condition was "Flag all bombs + uncover all non-bombs" not just the uncover all non-bombs
// Inputs: gameHandler Object to check board
// Outputs: checks to see if the game is in a win/lost state and edits that if needed
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
	if allNonBombsUncovered {
		handler.gameOver = true
		handler.win = true
	}
}

// Zhang: enabled AI functions (temp)
func (handler *Gamehandler) setAIEnabled(enabled bool) {
	handler.aiEnabled = enabled
	handler.aiTurn = false
}

func (handler *Gamehandler) setSolverEnabled(enabled bool) {
	handler.aiSolver = true
	handler.aiTurn = false
}

// Zhang: helper function for AI to take it move
func (handler *Gamehandler) RunAIMove() {
	if handler.gameOver {
		return
	}
	fmt.Println("aiDifficulty: ", handler.aiDifficulty)
	switch handler.aiDifficulty {
	case "Easy":
		if handler.aiSolver {
			fmt.Println("AI Solver making a move...")
			EasyAIMove(handler)
			time.Sleep(500 * time.Millisecond) // Pause for half a second between moves
		} else {
			EasyAIMove(handler)
		}
	case "Medium":
		if handler.aiSolver {
			for !handler.gameOver {
				fmt.Println("AI Solver making a move...")
				MediumAIMove(handler)
				time.Sleep(500 * time.Millisecond) // Pause for half a second between moves
			}
		} else {
			MediumAIMove(handler)
		}
	case "Hard":
		if handler.aiSolver {
			for !handler.gameOver {
				fmt.Println("AI Solver making a move...")
				HardAIMove(handler)
				time.Sleep(500 * time.Millisecond) // Pause for half a second between moves
			}
		} else {
			HardAIMove(handler)
		}
	}
}
