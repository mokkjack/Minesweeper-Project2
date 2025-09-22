package config

/*
Prologue

Authors: Adam Berry, Barrett Brown, Jonathan Gott, Alex Phibbs, Minh Vu
Creation Date: 9/11/2025

Description:
- This file hosts some basic config file bits to be called inside the go files instead of relying on magic numbers

Functions:
- None: Constant containing bunch of variables to be called


Inputs:
- None

Outputs:
- Int Value of the constant
*/

const (
	BoardSize    = 10 // Used to declare how many mines in the board
	MinMines     = 10 // Used to decide/display minimum allowed mines
	MaxMines     = 20 // Used to decide/display maximum allowed mines
	WindowHeight = 500 // Used to declare the window borders
	WindowWidth  = 500 // Used to declare window borders
	FixedWinSize = true // Bool to disallow adjusting window size
	GridSpacing  = (WindowHeight) / (BoardSize + 1) // Used to adjust how much space in the grid
)
