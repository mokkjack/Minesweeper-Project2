/*
Prologue

Authors: Adam Berry, Barrett Brown, Jonathan Gott, Alex Phibbs, Minh Vu
Creation Date: 9/21/2025

Description:
- This file implements the setup screen for the game, where the user is able to choose the number of
mines that they want in the game and validates that number if it's within the range. Afterwards it swaps the current view for the minesweeper view allowing the game to start

Functions:
- LoadSetupInfo: This loads the initial setup screen and asks the user for the number of mines.
Upon a valid entry, it'll create a new game and replaces the window with the game board.

Inputs:
- Mine count from the user

Outputs:
- Either the Minesweeper board, or the error message depending on if the user entered a mine number in range or not. (Error message will re-prompt for input)

*/

package components

import (
	"fmt"
	"strconv"

	"minesweeper/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// LoadSetupInto sets the setup form as content of an existing window.
// On submit it validates, builds the game UI, and replaces the window content.
// Inputs: the fyne window itself
// Outputs: Displays the window for the user
func LoadSetupInto(win fyne.Window) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(fmt.Sprintf("Enter mine count (%d-%d)", config.MinMines, config.MaxMines))
	entry.SetText(fmt.Sprintf("%d", config.MinMines))
	errLabel := widget.NewLabel("")

  // Create the "Setup window start"
	start := widget.NewButton("Start Game", func() {
    // Checks if entered value is int and not something random
		n, err := strconv.Atoi(entry.Text)
		if err != nil {
			errLabel.SetText("Please enter a valid integer.")
			return
		}
    // Bound checks
		maxAllowed := config.BoardSize*config.BoardSize - 1
		if n < config.MinMines || n > config.MaxMines {
			errLabel.SetText(fmt.Sprintf("Mine count must be between %d and %d.", config.MinMines, config.MaxMines))
			return
		}
		if n > maxAllowed {
			errLabel.SetText("Too many mines for this board size.")
			return
		}

		// Build the minesweeper game UI and swap it in to work
		h := NewGameHandler(n)
		board := GetBoard(&h)
		ui := SetupGameGraphics(board, &h)
		win.SetContent(ui)
	})

  // Creates a vertical box and shows it to display the setup to the user
	form := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Select number of mines (%d-%d):", config.MinMines, config.MaxMines)),
		entry,
		start,
		errLabel,
	)
	win.SetContent(container.NewPadded(form))
}
