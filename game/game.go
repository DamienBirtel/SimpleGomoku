package game

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	WINCON    = 4
	BOARDSIZE = 7
)

// Board represents the board on which the game takes place
type Board [BOARDSIZE][BOARDSIZE]int

// NewBoard retruns a new instance of a pointer to Board
func NewBoard() *Board {
	return &Board{}
}

func (b *Board) horizontalWin(x int, y int, stone int) bool {
	scoreToGet := WINCON - 1
	for i := x + 1; i < BOARDSIZE; i++ {
		if b[y][i] != stone {
			break
		}
		scoreToGet--
	}
	for i := x - 1; i >= 0; i-- {
		if b[y][i] != stone {
			break
		}
		scoreToGet--
	}
	if scoreToGet <= 0 {
		return true
	}
	return false
}

func (b *Board) verticalWin(x int, y int, stone int) bool {
	scoreToGet := WINCON - 1
	for i := y + 1; i < BOARDSIZE; i++ {
		if b[i][x] != stone {
			break
		}
		scoreToGet--
	}
	for i := y - 1; i >= 0; i-- {
		if b[i][x] != stone {
			break
		}
		scoreToGet--
	}
	if scoreToGet <= 0 {
		return true
	}
	return false
}

func (b *Board) firstDiagonalCheck(x int, y int, stone int) bool {
	scoreToGet := WINCON - 1
	i, j := x+1, y+1
	for i < BOARDSIZE && j < BOARDSIZE {
		if b[j][i] != stone {
			break
		}
		i, j = i+1, j+1
		scoreToGet--
	}
	i, j = x-1, y-1
	for i >= 0 && j >= 0 {
		if b[j][i] != stone {
			break
		}
		i, j = i-1, j-1
		scoreToGet--
	}
	if scoreToGet <= 0 {
		return true
	}
	return false
}

func (b *Board) secondDiagonalCheck(x int, y int, stone int) bool {
	scoreToGet := WINCON - 1
	i, j := x-1, y+1
	for i >= 0 && j < BOARDSIZE {
		if b[j][i] != stone {
			break
		}
		i, j = i-1, j+1
		scoreToGet--
	}
	i, j = x+1, y-1
	for i < BOARDSIZE && j >= 0 {
		if b[j][i] != stone {
			break
		}
		i, j = i+1, j-1
		scoreToGet--
	}
	if scoreToGet <= 0 {
		return true
	}
	return false
}

func (b *Board) diagonalWin(x int, y int, stone int) bool {
	win := b.firstDiagonalCheck(x, y, stone)
	if !win {
		win = b.secondDiagonalCheck(x, y, stone)
	}
	return win
}

// IsMoveWinning checks if the x and y placed stone is a winning move
func (b *Board) IsMoveWinning(x int, y int, stone int) bool {
	return (b.horizontalWin(x, y, stone) || b.verticalWin(x, y, stone) || b.diagonalWin(x, y, stone))
}

// Play places a new stone on the given coordinates and returns
// the winning player (1 or -1) if a player won, or 0 by default.
// It returns an error if a stone is already on these coordinates.
func (b *Board) Play(x int, y int, stone int) (int, error) {

	if b[y][x] != 0 {
		return 0, errors.New("ERROR: you can't place a stone on coordinates X:" + strconv.Itoa(x) + " Y:" + strconv.Itoa(y) + "\n       there is already a stone placed there\n")
	}
	b[y][x] = stone

	if b.IsMoveWinning(x, y, stone) {
		return stone, nil
	}

	return 0, nil
}

// Print prints the board in a pretty way with
// X symbols for player 1 and O symbols for player -1
func (b *Board) Print() {
	fmt.Println(" Y\\X |  0  |  1  |  2  |  3  |  4  |  5  |  6  |\n")
	for i, line := range b {
		fmt.Print(" ", i, "  |")
		for _, cell := range line {
			switch cell {
			case 1:
				fmt.Print("   X  ")
			case -1:
				fmt.Print("   O  ")
			default:
				fmt.Print("   .  ")
			}
		}
		fmt.Println("|\n")
	}
}
