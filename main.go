package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DamienBirtel/SimpleGomoku/game"
	"github.com/DamienBirtel/SimpleGomoku/mcts"
)

func main() {
	var x, y int
	var err error

	scanner := bufio.NewScanner(os.Stdin)
	board := game.NewBoard()
	tree := mcts.NewTree(board, 1)

	stopChan := make(chan struct{})
	playChan := make(chan struct{})
	c := make(chan [2]int)
	go tree.ComputeMCTS(c, stopChan, playChan)

	winner := 0
	round := 0
	for winner == 0 {
		if err == nil {
			board.Print()
		}
		switch round % 2 {
		case 1:
			playChan <- struct{}{}
			move := <-c
			_, _ = board.Play(move[0], move[1], -1)
			round++
		default:
			fmt.Print("Your move (x y): ")
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &x, &y)
			c <- [2]int{x, y}
			winner, err = board.Play(x, y, 1)

			if err != nil {
				fmt.Println(err)
				break
			}
			round++
		}
	}
	board.Print()
	fmt.Println("GG ", winner)
}
