# EECS 581 Project 1

## Getting Started:

- [Fyne Getting Started](https://docs.fyne.io/started/)
- [Setting up Golang](https://go.dev/doc/tutorial/getting-started)
- [Guides on setting up Go modules/Using Go if you haven't](https://go.dev/doc/tutorial/create-module)

## General Layout/Code execution:

- Config options i.e: Min/max mines and what not are pre-defined in the constants.go file
- All execution starts in "main.go" this is started by running make or go run .
- Afterwards main.go will contact setup.go to create a window and ask the user how many mines they want
- Upon declaring how many mines will be "in play" it will connect to ui-handler/game-handler.go
- main.go: General entry point for the user, in here it will call to setup.go to "show" the initial window then swap view in that window to the minesweeper game
- ui-handler.go is used to display the cells with the neighbor numbers/state/grab initial left/right click (uncover/flag) and do what needs to be done there
  - Set up cells/grid
  - Grab clicks/"push" clicked row/col onto other func in game-handler.go
  - updateCellTexts: Updates text of all cells visually if needed (Necessary for when the "regeneration" happens if first click is bomb since if not refreshed it will still display it as a bomb)
  - applyOverlayStates: Used to "refresh" the state of the pre-placed cells based on updates from flood/other actions
  - SetupGameGraphics: Used to generate initial cells/create win & loss button (Sets invisible at start so later when edited it can "show")
  - updateGameUI: Used as a general "Update all states" flow, allows you to update text/visual states then after checks if the win/lost condition needs to show, if so show them
- game-handler.go handles most of the "game logic" rules, this is used to adjust some 2D-Arrays that the UI handler looks out to figure out "what to display"
  - Initial Game setup/bomb placement
  - Neighbor Counting
  - Click function is used to do a couple things including:
    - Move bomb/regenerate if first click = bomb
    - "Flood reveal" on click
    - Handle if clicked on bomb
    - Check win condition
  - Flagging on 2D-array
