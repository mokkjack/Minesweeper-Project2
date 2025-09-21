package main

import (
	"minesweeper/components"
	"minesweeper/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

//var numberOfMines int = 10   // User Determined, can be 10 or 20

func main() {
	a := app.New()
	window := a.NewWindow("Minesweeper")
	window.Resize(fyne.NewSize(config.WindowHeight, config.WindowWidth))
	window.SetFixedSize(config.FixedWinSize)

	gameHandler := components.NewGameHandler()
	board := components.GetBoard(&gameHandler)
	gameWindow := components.SetupGameGraphics(board, &gameHandler) //had to change the uiboard function so this needed the use the board variable instead of getboard again
	window.SetContent(gameWindow)

	// Fake functions lie in "ui-handler.go" and are meant just to simulate an uncover (directly calls RevealZero)/flag operation
	components.FakeUncover(&gameHandler, 2, 3)
	components.FakeFlag(board, 0, 0)
	window.ShowAndRun()
}
