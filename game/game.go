package game

// Board represents the board on which the game takes place
type Board [5][5]int

func NewBoard() *Board {
	return &Board{}
}
