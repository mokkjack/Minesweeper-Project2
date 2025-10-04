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
	"image/color"
	"minesweeper/config"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// LoadSetupInto sets the setup form as content of an existing window.
// On submit it validates, builds the game UI, and replaces the window content.
// Inputs: the fyne window itself
// Outputs: Displays the window for the user
func LoadSetupInto(win fyne.Window) {
	//Title Card
	title := canvas.NewText("MINESWEEPER 2", color.RGBA{0, 255, 0, 255})
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 32
	titlePlace := container.NewCenter(title)

	//Play Button
	playButton := widget.NewButton("Play Game", func() {
		gameSelect(win)
	})

	//Exit Button
	exitButton := widget.NewButton("Exit", func() {
		win.Close()
	})

	from := container.NewVBox(
		titlePlace,
		playButton,
		exitButton,
	)

	win.SetContent(container.NewPadded(from))
}

// Game Select Screen
func gameSelect(win fyne.Window) {
	modelLabel := widget.NewLabel("Choose Game Mode:")

	singleButton := widget.NewButton("Single Player", func() {
		showMineSetup(win, "Single", "Play")
	})
	aiButton := widget.NewButton("AI 1v1 Mode", func() {
		showAImode(win, "comp")
	})
	solverButton := widget.NewButton("AI Solver Mode", func() {
		showAImode(win, "Solver")
	})

	from := container.NewVBox(
		modelLabel,
		singleButton,
		aiButton,
		solverButton,
	)
	win.SetContent(container.NewPadded(from))
}

// AI Mode Screen || Zhang: show AI mode setup
func showAImode(win fyne.Window, mode string) {
	label := widget.NewLabel("Select AI Difficulty:")
	easy := widget.NewButton("Easy", func() {
		if mode == "comp" {
			showMineSetup(win, "AI", "Easy")
		} else {
			showMineSetup(win, "Solve", "Easy")
		}
	})
	medium := widget.NewButton("Medium", func() {
		if mode == "comp" {
			showMineSetup(win, "AI", "Medium")
		} else {
			showMineSetup(win, "Solve", "Medium")
		}
	})
	hard := widget.NewButton("Hard", func() {
		if mode == "comp" {
			showMineSetup(win, "AI", "Hard")
		} else {
			showMineSetup(win, "Solve", "Hard")
		}
	})
	from := container.NewVBox(label, easy, medium, hard)
	win.SetContent(container.NewPadded(from))
}

// Mine Setup Screen
func showMineSetup(win fyne.Window, mode string, option string) {
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
		h := NewGameHandler(n)
		//Zhang: Apply selected mode
		fmt.Print("Selected mode: ", mode, " with option: ", option, "\n")
		if mode == "AI" {
			h.setAIEnabled(true)
			h.aiDifficulty = option
		} else if mode == "Solve" {
			fmt.Println("Single Player - Solve mode")
			h.setSolverEnabled(true)
			h.aiDifficulty = option
		}
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
