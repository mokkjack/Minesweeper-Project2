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
func LoadSetupInto(win fyne.Window) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(fmt.Sprintf("Enter mine count (%d-%d)", config.MinMines, config.MaxMines))
	entry.SetText("10")
	errLabel := widget.NewLabel("")

	start := widget.NewButton("Start Game", func() {
		n, err := strconv.Atoi(entry.Text)
		if err != nil {
			errLabel.SetText("Please enter a valid integer.")
			return
		}
		maxAllowed := config.BoardSize*config.BoardSize - 1
		if n < config.MinMines || n > config.MaxMines {
			errLabel.SetText(fmt.Sprintf("Mine count must be between %d and %d.", config.MinMines, config.MaxMines))
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
		widget.NewLabel(fmt.Sprintf("Select number of mines (%d-%d):", config.MinMines, config.MaxMines)),
		entry,
		start,
		errLabel,
	)
	win.SetContent(container.NewPadded(form))
}
