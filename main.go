package main

import (
	"minesweeper/config"
	"minesweeper/components"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)
//var numberOfMines int = 10   // User Determined, can be 10 or 20


func main() {
	a := app.New()
	window := a.NewWindow("Minesweeper")
	window.Resize(fyne.NewSize(config.WindowHeight,config.WindowWidth))
	window.SetFixedSize(config.FixedWinSize)
	game := components.SetupGame()
	window.SetContent(game)
	window.ShowAndRun()
}
