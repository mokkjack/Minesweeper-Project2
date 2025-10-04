/*
Prologue

Authors: Adam Berry, Barrett Brown, Jonathan Gott, Alex Phibbs, Minh Vu
Creation Date: 9/11/2025

Description:
- This file handles the GUI for a Minesweeper game. We're mostly using the Fyne GUI library for this. It uses the backend stuff from
game-handler and creates the GUI from that. This will create the visual grid, create the text in each cell based on the state (covered, uncovered, or flags) and the underlying text.
It also displays the win/lose message.

Functions:
- SetupGameGraphics: Initializes all GUI parts for the board creating the initial cells/win & lose message (keeping them inivisble)

- Tapped: Handles all left clicks

- TappedSecondary: Handles all right clicks

- applyOverlayStates: Update overlay visibility and colors based on the state of the cell (uncovered, covered, flagged, etc.). This is meant so when updating the states of the cells upon clicking it will properly reflect it on the visual side

- updateCellTexts: Updates the text inside the cell, useful for if the board had to be regenerated due to a "first left click on bomb" as the numbers in the 2d array wouldn't be updated alone by applyOverlayStates

- UpdateGameUI: Used to refresh both celltext/overlay states in the correct order as well as check the win condition upon which it will show some text overlays (i.e. end of game message)

Input:
- Board state from game-handler
- Player mouse clicks

Output:
- GUI visual update to the board i.e: uncover/flag
- End of game messages

*/

package components

import (
	"minesweeper/config"
	"time"

	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	cellOverlays [][]*canvas.Rectangle
	cellFlags    [][]*canvas.Text
	cellTexts    [][]*canvas.Text
	gameMsg      *canvas.Text

	gameOverContainer *fyne.Container
	newGameButton     *widget.Button
	titleScreenButton *widget.Button
)

type clickableRect struct {
	*canvas.Rectangle
	row     int
	col     int
	handler *Gamehandler
}

var _ fyne.Tappable = (*clickableRect)(nil)
var _ fyne.SecondaryTappable = (*clickableRect)(nil)

/*
Called upon left click, will check if game is already over (Not allow gameplay past loss/win) and then afterwards calls game-handler.go's Click function to handle the backend click and then updates the game ui based on what that did
*/
func (c *clickableRect) Tapped(_ *fyne.PointEvent) {
	if c.handler.gameOver {
		return
	}
	// Zhang: prevent user from clicking when it's AI's turn
	if c.handler.aiEnabled && c.handler.aiTurn {
		return
	}
	sq := &c.handler.board[c.row][c.col]

	if sq.state == Uncovered || sq.state == Flagged {
		return
	}
	c.handler.Click(c.row, c.col)
	UpdateGameUI(c.handler)

	if c.handler.aiEnabled && !c.handler.gameOver {
		c.handler.aiTurn = true
		c.handler.RunAIMove()
		c.handler.aiTurn = false
		UpdateGameUI(c.handler)
	} else if c.handler.aiSolver && !c.handler.gameOver {
		c.handler.aiTurn = true
		go func() { // Run the AI solver in a separate goroutine
			for !c.handler.gameOver {
				c.handler.RunAIMove()
				UpdateGameUI(c.handler)
				time.Sleep(1000 * time.Millisecond) // Pause for one second between moves
			}
			c.handler.aiTurn = false
		}()
	}
}

/*
Called upon right click, checks if game over and then turns the underlining 2d-array to have a flag state and then refresh the game ui
*/
func (c *clickableRect) TappedSecondary(_ *fyne.PointEvent) {
	if c.handler.gameOver { // ignore flags after game over
		return
	}
	if c.handler.aiEnabled && c.handler.aiTurn { // Zhang: prevent user from flagging when it's AI's turn
		return
	}
	sq := &c.handler.board[c.row][c.col]
	if sq.state == Uncovered {
		return
	}
	c.handler.ToggleFlag(c.row, c.col)
	UpdateGameUI(c.handler)

	if c.handler.aiEnabled && !c.handler.gameOver { // Zhang: let AI make a move after user right clicks
		c.handler.aiTurn = true
		EasyAIMove(c.handler)
		c.handler.aiTurn = false
		UpdateGameUI(c.handler)
	}
}

// Helper function: Used to simplify r.move() operations
func cellPos(col, row int) fyne.Position {
	return fyne.NewPos(
		float32(col*config.GridSpacing),
		float32(row*config.GridSpacing),
	)
}

// This Function is Intended to be used as a one time initializer for the game's UI components
// Inputs: 2D-Array of the board and the gameHandler object to get the context of the object for the click handler
// Outputs: A fyne container which can store multiple elements
func SetupGameGraphics(board [][]Square, handler *Gamehandler) *fyne.Container {
	var columnNames string = "abcdefghijklmnopqrstuvwxyz"

	// Initialize storage variables for Overlays/flags/Textboxes
	// Create "Cells" on top of each box to show/not show depending on state
	cellOverlays = make([][]*canvas.Rectangle, config.BoardSize)
	cellFlags = make([][]*canvas.Text, config.BoardSize)
	for r := range cellOverlays {
		cellOverlays[r] = make([]*canvas.Rectangle, config.BoardSize)
		cellFlags[r] = make([]*canvas.Text, config.BoardSize)
	}
	cellTexts = make([][]*canvas.Text, config.BoardSize)
	for r := range cellTexts {
		cellTexts[r] = make([]*canvas.Text, config.BoardSize)
	}

	// Used as an array to "loop" over in order so to ensure proper "layering" of each item (did * 5 just to ensure extra space not really needed to be this big)
	objects := make([]fyne.CanvasObject, 0, (config.BoardSize+1)*(config.BoardSize+1)*5)

	// Loop overboard setting row/column headers (row/col == 0 lines), if the cell is a body cell instead we "draw" the text for that cell (Bomb/neighbors)
	for row := 0; row < (config.BoardSize + 1); row++ {
		for col := 0; col < (config.BoardSize + 1); col++ {
			if row == 0 && col == 0 {
				continue
			} else if row == 0 {
				r := canvas.NewText(columnNames[col-1:col], color.RGBA{255, 255, 255, 255})
				r.TextSize = float32(config.GridSpacing) / 2
				sz := r.MinSize()
				cell := float32(config.GridSpacing)
				x := float32(col*config.GridSpacing) + (cell-sz.Width)/2
				y := float32(0*config.GridSpacing) + (cell-sz.Height)/2
				r.Move(fyne.NewPos(x, y))
				objects = append(objects, r)
			} else if col == 0 {
				r := canvas.NewText(strconv.Itoa(row), color.RGBA{255, 255, 255, 255})
				r.TextSize = float32(config.GridSpacing) / 2
				sz := r.MinSize()
				cell := float32(config.GridSpacing)
				x := float32(0*config.GridSpacing) + (cell-sz.Width)/2
				y := float32(row*config.GridSpacing) + (cell-sz.Height)/2
				r.Move(fyne.NewPos(x, y))
				objects = append(objects, r)
			} else {
				// Draw underlying cell content (bomb or number)
				c := board[row-1][col-1]
				var txt string
				if c.isBomb {
					txt = "b"
				} else if c.numValue != 0 {
					txt = strconv.Itoa(c.numValue)
				}
				base := canvas.NewText(txt, color.RGBA{0, 255, 0, 255})
				base.TextSize = config.GridSpacing / 2

				// Center Text
				size := base.MinSize()
				cellSize := float32(config.GridSpacing)

				x := float32(col*config.GridSpacing) + (cellSize-size.Width)/2
				y := float32(row*config.GridSpacing) + (cellSize-size.Height)/2
				base.Move(fyne.NewPos(x, y))

				objects = append(objects, base)
				cellTexts[row-1][col-1] = base
			}
		}
	}

	// As of this point the "cells" above havce the underlining neighbor/bomb/row & col header but the covering "cell" bit that you can click isn't on there so this re loops through and places them
	// We first create the rectangles objects and place them where they go setting their colors and what not
	for rw := 0; rw < config.BoardSize; rw++ {
		for c := 0; c < config.BoardSize; c++ {
			// overlay rectangle
			overlay := canvas.NewRectangle(color.NRGBA{R: 60, G: 60, B: 60, A: 255})
			overlay.Resize(fyne.NewSize(float32(config.GridSpacing), float32(config.GridSpacing)))
			overlay.Move(cellPos(c+1, rw+1)) // adjust for header

			overlay.StrokeColor = color.NRGBA{R: 30, G: 30, B: 30, A: 255}
			overlay.StrokeWidth = 1

			clickable := &clickableRect{
				Rectangle: overlay,
				row:       rw,
				col:       c,
				handler:   handler,
			}

			cellOverlays[rw][c] = overlay

			// Flag
			flag := canvas.NewText("F", color.NRGBA{R: 220, G: 40, B: 40, A: 255})
			flag.TextSize = config.GridSpacing / 2
			flag.TextStyle.Bold = true
			size := flag.MinSize()

			// Center Text
			cellSize := float32(config.GridSpacing)
			x := float32((c+1)*config.GridSpacing) + (cellSize-size.Width)/2
			y := float32((rw+1)*config.GridSpacing) + (cellSize-size.Height)/2
			flag.Move(fyne.NewPos(x, y))
			cellFlags[rw][c] = flag

			objects = append(objects, clickable.Rectangle, clickable, flag)
		}
	}

	// Finally we create the "end game" message object
	// We also make sure it is centered (hidden initially)
	gameMsg = canvas.NewText("", color.White)
	gameMsg.TextStyle.Bold = true
	gameMsg.TextSize = float32(config.GridSpacing) * 0.9

	newGameButton = widget.NewButton("Restart", func() {
		win := fyne.CurrentApp().Driver().AllWindows()[0]
		mineCount := handler.totalMines
		h := NewGameHandler(mineCount)
		if handler.aiEnabled {
			h.setAIEnabled(true)
			h.aiDifficulty = handler.aiDifficulty
		}
		board := GetBoard(&h)
		ui := SetupGameGraphics(board, &h)
		win.SetContent(ui)
	})

	titleScreenButton = widget.NewButton("Title Screen", func() {
		win := fyne.CurrentApp().Driver().AllWindows()[0]
		LoadSetupInto(win)
	})

	gameOverContainer = container.NewVBox(
		gameMsg,
		container.NewHBox(
			newGameButton,
			titleScreenButton,
		),
	)
	// center over the whole board (headers + grid)
	totalPx := float32(config.GridSpacing * (config.BoardSize + 1))
	ms := gameOverContainer.MinSize()
	gameOverContainer.Move(fyne.NewPos((totalPx-ms.Width)/2, (totalPx-ms.Height)/2))
	gameOverContainer.Hide()

	objects = append(objects, gameOverContainer)

	// Call to apply overlay states as now that the object itself is "fleshed out" we can actually display it
	applyOverlayStates(board)

	return container.NewWithoutLayout(objects...)
}

/*
Used as a cell refresher, all this does is go back over all the cells in the cell and check state and apply correct cell properties according to it
Inputs: 2D-Array of the boards cells
Outputs: None, just refreshing the underlying values
*/
func applyOverlayStates(board [][]Square) {
	for r := 0; r < config.BoardSize; r++ {
		for c := 0; c < config.BoardSize; c++ {
			ov := cellOverlays[r][c]
			fl := cellFlags[r][c]
			switch board[r][c].state {
			case Covered:
				ov.FillColor = color.NRGBA{R: 60, G: 60, B: 60, A: 255}
				ov.Refresh()
				fl.Hide()
			case Uncovered:
				ov.FillColor = color.NRGBA{R: 60, G: 60, B: 60, A: 0}
				ov.Refresh()
				fl.Hide()
			case Flagged:
				ov.FillColor = color.NRGBA{R: 60, G: 60, B: 60, A: 255}
				ov.Refresh()
				fl.Show()
			}
		}
	}
}

/*
Used mainly as a verifier upon the case that cell text values changed, this specifically happens if needing to regenerate the board due to your first click being a bomb
Inputs: 2D-Array of the boards cells
Outputs: None, just refreshing the underlying values
*/
func updateCellTexts(board [][]Square) {
	for r := 0; r < config.BoardSize; r++ {
		for c := 0; c < config.BoardSize; c++ {
			t := cellTexts[r][c]
			if t == nil {
				continue
			}

			// decide what to show
			var txt string
			if board[r][c].isBomb {
				txt = "b"
			} else if board[r][c].numValue != 0 {
				txt = strconv.Itoa(board[r][c].numValue)
			} else {
				txt = "" // empty for zeros
			}

			// update text and keep it centered
			if t.Text != txt {
				t.Text = txt
				t.TextSize = config.GridSpacing / 2
				sz := t.MinSize()
				cell := float32(config.GridSpacing)
				// +1,+1 because the board is offset by headers
				x := float32((c+1)*config.GridSpacing) + (cell-sz.Width)/2
				y := float32((r+1)*config.GridSpacing) + (cell-sz.Height)/2
				t.Move(fyne.NewPos(x, y))
				t.Refresh()
			}
			if board[r][c].markedByAI {
				t.Color = color.RGBA{255, 255, 0, 255} // Yellow
			}
			t.Refresh()
		}
	}
}

/*
Inputs: Game handler object for the context
Outputs: None, just refreshes UI/Shows win condition to screen
*/
func UpdateGameUI(h *Gamehandler) {
	updateCellTexts(h.board)
	applyOverlayStates(h.board)
	if h.gameOver { //play again + title button
		if h.win {
			gameMsg.Text = "You Win!"
			gameMsg.TextStyle.Bold = true
			gameMsg.Color = color.RGBA{R: 255, G: 222, B: 33, A: 255}
		} else {
			gameMsg.Text = "Game Over"
			gameMsg.TextStyle.Bold = true
			gameMsg.Color = color.RGBA{R: 220, A: 255}
		}
		gameMsg.Refresh()

		totalPx := float32(config.GridSpacing * (config.BoardSize + 1))
		ms := gameOverContainer.MinSize()
		gameOverContainer.Move(fyne.NewPos((totalPx-ms.Width)/2, (totalPx-ms.Height)/2))

		gameOverContainer.Show()
		gameOverContainer.Refresh()
	} else {
		gameOverContainer.Hide()
	}
}
