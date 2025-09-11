package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func main() {
	a := app.New()
	w := a.NewWindow("Grid Layout")
  grid := container.New(layout.NewGridLayout(11))

  for row := 0; row < 11; row++{
    for col := 0; col < 11; col++{
      if row == 0 && col == 0 {

      } else if row == 0{

      } else if col == 0{

      } else{
        
      }
    }
  }

	w.SetContent(grid)
	w.ShowAndRun()
}
