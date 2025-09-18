package components

import (
	"minesweeper/config"
)

// mine generator functions needs the location of the first click, and the number of mines
func (h *gamehandler) mine_generator(first_click_row, first_click_col, num_mines int) {
	// find the cell that was the first click
	first_clicked_cell := first_click_row*config.BoardSize + first_click_col

	// represents the total number of cells
	num_cells := config.BoardSize * config.BoardSize

	// this slice will have all locations where mines CAN go
	possible_mine_locations := make([]int, 0, num_cells)

	// this for-loop finds every cell that is not the first clicked cell
	// and adds it to the list of possible mine locations
	for row := 0; row < config.BoardSize; row++ {
		for col := 0; col <= config.BoardSize; col++ {
			// find current cell
			cell_id := row*config.BoardSize + col

			// If we come upon the first cell that was clicked, skip it so
			// that it does not become a possible mine location
			if cell_id != first_clicked_cell {
				possible_mine_locations = append(possible_mine_locations, cell_id)
			}
		}
	}

	// Shuffle the possible mine locations randomly
	h.rng.Shuffle(len(possible_mine_locations), func(i int, j int) {
		possible_mine_locations[i], possible_mine_locations[j] = possible_mine_locations[j], possible_mine_locations[i]
	})

	// Add mines to the first 10 or 20 cells in the shuffled list of possible mine locations
	for i := 0; i < num_mines; i++ {
		// get the cell location
		cell_id := possible_mine_locations[i]

		// convert cell id to row, column
		row := cell_id / config.BoardSize
		col := cell_id % config.BoardSize

		// add a mine to the cell
		h.board[row][col].isBomb = true
	}

}
