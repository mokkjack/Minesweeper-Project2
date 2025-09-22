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
  cellTexts    [][]*canvas.Text
  gameMsg *canvas.Text
)

type clickableRect struct {
	*canvas.Rectangle
	row        int
	col        int
	handler    *Gamehandler
}

var _ fyne.Tappable = (*clickableRect)(nil)
var _ fyne.SecondaryTappable = (*clickableRect)(nil)

func (c *clickableRect) Tapped(_ *fyne.PointEvent) {
	if c.handler.gameOver {
    return
  }
	c.handler.Click(c.row, c.col)
	updateGameUI(c.handler)
}

func (c *clickableRect) TappedSecondary(_ *fyne.PointEvent) {
	if c.handler.gameOver { // ignore flags after game over
    return
  }
  c.handler.ToggleFlag(c.row, c.col)
	updateGameUI(c.handler)
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
  cellTexts = make([][]*canvas.Text, config.BoardSize)
  for r := range cellTexts {
    cellTexts[r] = make([]*canvas.Text, config.BoardSize)
  }


	objects := make([]fyne.CanvasObject, 0, (config.BoardSize+1)*(config.BoardSize+1)*5)

	for row := 0; row < (config.BoardSize + 1); row++ {
		for col := 0; col < (config.BoardSize + 1); col++ {
			if row == 0 && col == 0 {
				continue
			} else if row == 0 {
				r := canvas.NewText(columnNames[col-1:col], color.White)
				r.TextSize = float32(config.GridSpacing) / 2
        sz := r.MinSize()
        cell := float32(config.GridSpacing)
        x := float32(col*config.GridSpacing) + (cell - sz.Width)/2
        y := float32(0*config.GridSpacing)   + (cell - sz.Height)/2
        r.Move(fyne.NewPos(x, y))
				colHeaders = append(colHeaders, r)
				objects = append(objects, r)
			} else if col == 0 {
				r := canvas.NewText(strconv.Itoa(row), color.White)
				r.TextSize = float32(config.GridSpacing) / 2
        sz := r.MinSize()
        cell := float32(config.GridSpacing)
        x := float32(0*config.GridSpacing)   + (cell - sz.Width)/2
        y := float32(row*config.GridSpacing) + (cell - sz.Height)/2
        r.Move(fyne.NewPos(x, y))
				rowHeaders = append(rowHeaders, r)
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
				base := canvas.NewText(txt, color.White)
				base.TextSize = config.GridSpacing / 2

				// Center Text
				size := base.MinSize()
				cellSize := float32(config.GridSpacing)

				x := float32(col*config.GridSpacing) + (cellSize-size.Width) / 2
				y := float32(row*config.GridSpacing) + (cellSize-size.Height) / 2
				base.Move(fyne.NewPos(x, y))

				objects = append(objects, base)
        cellTexts[row-1][col-1] = base
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

  // Centered end-game message (hidden initially)
  gameMsg = canvas.NewText("", color.White)
  gameMsg.TextStyle.Bold = true
  gameMsg.TextSize = float32(config.GridSpacing) * 0.9

  // center over the whole board (headers + grid)
  totalPx := float32(config.GridSpacing * (config.BoardSize + 1))
  ms := gameMsg.MinSize()
  gameMsg.Move(fyne.NewPos((totalPx-ms.Width)/2, (totalPx-ms.Height)/2))
  gameMsg.Hide()

  objects = append(objects, gameMsg)

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

func updateCellTexts(board [][]Square) {
  for r := 0; r < config.BoardSize; r++ {
    for c := 0; c < config.BoardSize; c++ {
      t := cellTexts[r][c]
      if t == nil { continue }

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
        x := float32((c+1)*config.GridSpacing) + (cell - sz.Width)/2
        y := float32((r+1)*config.GridSpacing) + (cell - sz.Height)/2
        t.Move(fyne.NewPos(x, y))
        t.Refresh()
      }
    }
  }
}


func updateGameUI(h *Gamehandler) {
  updateCellTexts(h.board)
  applyOverlayStates(h.board)

  if h.gameOver {
    if h.win {
      gameMsg.Text = "You Win!"
      gameMsg.Color = color.RGBA{G: 220, A: 255}
    } else {
      gameMsg.Text = "Game Over"
      gameMsg.Color = color.RGBA{R: 220, A: 255}
    }

    totalPx := float32(config.GridSpacing * (config.BoardSize + 1))
    ms := gameMsg.MinSize()
    gameMsg.Move(fyne.NewPos((totalPx-ms.Width)/2, (totalPx-ms.Height)/2))

    
    gameMsg.Show()
    gameMsg.Refresh()
  } else {
    gameMsg.Hide()
  }
}
