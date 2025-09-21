package components

import (
	"minesweeper/config"

	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

var (
	colHeaders []fyne.CanvasObject
	rowHeaders []fyne.CanvasObject

	cellOverlays [][]*canvas.Rectangle
	cellFlags    [][]*canvas.Text
)

type clickableRect struct {
	*canvas.Rectangle
	row        int
	col        int
	handler    *Gamehandler
	firstClick bool
}

var _ fyne.Tappable = (*clickableRect)(nil)
var _ fyne.SecondaryTappable = (*clickableRect)(nil)

func (c *clickableRect) Tapped(_ *fyne.PointEvent) {
	if c.firstClick {
		c.firstClick = false
	}
	c.handler.RevealZero(c.row, c.col)
	applyOverlayStates(c.handler.board)
}

func (c *clickableRect) TappedSecondary(_ *fyne.PointEvent) {
	sq := &c.handler.board[c.row][c.col]
	if sq.state == Flagged {
		sq.state = Covered
	} else if sq.state == Covered {
		sq.state = Flagged
	}
	applyOverlayStates(c.handler.board)
}

// Used to simplify r.move() operations
func cellPos(col, row int) fyne.Position {
	return fyne.NewPos(
		float32(col*config.GridSpacing),
		float32(row*config.GridSpacing),
	)
}

// This Function is Intended to be used as a one time initializer for the game's UI components
// Inputs: None
// Outputs: A fyne container which can store multiple elements
func SetupGameGraphics(board [][]Square, handler *Gamehandler) *fyne.Container { //added the game handler to this so it could use the logic from game handler.go Alex
	var columnNames string = "abcdefghijklmnopqrstuvwxyz"

	// Create "Cells" on top of each box to show/not show depending on state
	cellOverlays = make([][]*canvas.Rectangle, config.BoardSize)
	cellFlags = make([][]*canvas.Text, config.BoardSize)
	for r := range cellOverlays {
		cellOverlays[r] = make([]*canvas.Rectangle, config.BoardSize)
		cellFlags[r] = make([]*canvas.Text, config.BoardSize)
	}

	objects := make([]fyne.CanvasObject, 0, (config.BoardSize+1)*(config.BoardSize+1)*5)

	for row := 0; row < (config.BoardSize + 1); row++ {
		for col := 0; col < (config.BoardSize + 1); col++ {
			if row == 0 && col == 0 {
				continue
			} else if row == 0 {
				r := canvas.NewText(columnNames[col-1:col], color.White)
				r.Move(cellPos(col, 0))
				colHeaders = append(colHeaders, r)
				objects = append(objects, r)
			} else if col == 0 {
				r := canvas.NewText(strconv.Itoa(row), color.White)
				r.Move(cellPos(0, row))
				rowHeaders = append(rowHeaders, r)
				objects = append(objects, r)
			} else {
				// Draw underlying cell content (bomb or number)
				c := board[row-1][col-1]
				var txt string
				if c.isBomb {
					txt = "b"
				} else {
					txt = strconv.Itoa(c.numValue)
				}
				base := canvas.NewText(txt, color.White)
				base.TextSize = config.GridSpacing / 2

				// Center Text
				size := base.MinSize()
				cellSize := float32(config.GridSpacing)

				x := float32(col*config.GridSpacing) + (cellSize-size.Width)/2
				y := float32(row*config.GridSpacing) + (cellSize-size.Height)/2
				base.Move(fyne.NewPos(x, y))

				objects = append(objects, base)
			}
		}
	}

	for rw := 0; rw < config.BoardSize; rw++ {
		for c := 0; c < config.BoardSize; c++ {
			// overlay rectangle
			overlay := canvas.NewRectangle(color.NRGBA{R: 60, G: 60, B: 60, A: 255})
			overlay.Resize(fyne.NewSize(float32(config.GridSpacing), float32(config.GridSpacing)))
			overlay.Move(cellPos(c+1, rw+1)) // adjust for header

			overlay.StrokeColor = color.NRGBA{R: 30, G: 30, B: 30, A: 255}
			overlay.StrokeWidth = 1

			clickable := &clickableRect{
				Rectangle:  overlay,
				row:        rw,
				col:        c,
				handler:    handler,
				firstClick: true,
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

			objects = append(objects, clickable, flag)
		}
	}

	applyOverlayStates(board)

	return container.NewWithoutLayout(objects...)
}

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

// Created to simulate an "Uncover operation"
// Using gamehandler object to call "board-manager.go"'s revealZero func without causing import cycle
func FakeUncover(h *Gamehandler, row, col int) {
	if h == nil || row < 0 || col < 0 || row >= len(h.board) || col >= len(h.board[row]) {
		return
	}
	h.RevealZero(row, col)
	applyOverlayStates(h.board) // repaint states
}

// Created to simulate a flag operation
func FakeFlag(board [][]Square, row, col int) {
	if row < 0 || col < 0 || row >= len(board) || col >= len(board[row]) {
		return
	}
	board[row][col].state = Flagged
	applyOverlayStates(board)
}
