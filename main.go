/*
Authors: Adam Berry, Barrett Brown, Jonathan Gott, Alex Phibbs, Minh Vu
Creation Date: 9/11/2025

Description: Initializes everything and especially the Fyne app. Sets up the main window and loads
the setup screen.
*/

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

	components.LoadSetupInto(window)
	window.ShowAndRun()
}
