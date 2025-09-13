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
)

// This Function is Intended to be used as a one time initializer for the game's UI components
// We'll need to modify this so that it can accept the gamestate as an argument 
// Inputs: None
// Outputs: A fyne container which can store multiple elements
func SetupGameGraphics(board [][]Square) *fyne.Container{
	var columnNames string = "abcdefghijklmnopqrstuvwxyz"

	var tempMines [10][10]int
	tempMines[0][2] = 1
	tempMines[4][8] = 1


	for row := 0; row < (config.BoardSize + 1); row++{
		for col := 0; col < (config.BoardSize + 1); col++{
			if row == 0 && col == 0 {
				continue
			} else if row == 0{
				r := canvas.NewText(columnNames[col-1:col], color.White)
				r.Move(fyne.NewPos(float32(col*config.GridSpacing), 0) )
				colHeaders = append(colHeaders, r)
			} else if col == 0{
				r := canvas.NewText(strconv.Itoa(row), color.White)
				r.Move(fyne.NewPos(0, float32(row*config.GridSpacing) ) )
				rowHeaders = append(rowHeaders, r)
			} else if board[row-1][col-1].isBomb{
				r := canvas.NewText("b", color.White)
				r.Move(fyne.NewPos(float32(col*config.GridSpacing), float32(row*config.GridSpacing) ) )
				rowHeaders = append(rowHeaders, r)
			} else {
				r := canvas.NewText(strconv.Itoa(board[row-1][col-1].numValue), color.White)
				r.Move(fyne.NewPos(float32(col*config.GridSpacing), float32(row*config.GridSpacing) ) )
				rowHeaders = append(rowHeaders, r)
			}
		}
	}
	all_items := append(rowHeaders, colHeaders...)
	layout := container.New(nil, all_items...)

	return layout 
}


