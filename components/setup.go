package components

import (
	"strconv"

	"minesweeper/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// LoadSetupInto sets the setup form as content of an existing window.
// On submit it validates, builds the game UI, and replaces the window content.
func LoadSetupInto(win fyne.Window) {
	const minMines = 10
	const maxMines = 20

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter mine count (10-20)")
	entry.SetText("10")
	errLabel := widget.NewLabel("")

	start := widget.NewButton("Start Game", func() {
		n, err := strconv.Atoi(entry.Text)
		if err != nil {
			errLabel.SetText("Please enter a valid integer.")
			return
		}
		maxAllowed := config.BoardSize*config.BoardSize - 1
		if n < minMines || n > maxMines {
			errLabel.SetText("Mine count must be between 10 and 20.")
			return
		}
		if n > maxAllowed {
			errLabel.SetText("Too many mines for this board size.")
			return
		}

		// Build the game UI and swap it in
		h := NewGameHandler(n)
		board := GetBoard(&h)
		ui := SetupGameGraphics(board, &h)
		win.SetContent(ui)
	})


	form := container.NewVBox(
		widget.NewLabel("Select number of mines (10-20):"),
		entry,
		start,
		errLabel,
	)
	win.SetContent(container.NewPadded(form))
}
