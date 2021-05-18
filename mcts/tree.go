package mcts

import (
	"math/rand"
	"time"

	"github.com/DamienBirtel/SimpleGomoku/game"
)

type Tree struct {
	root *Node
}

type Node struct {
	id            [2]int
	nbVisit       int
	nbWins        int
	playing       int
	state         game.Board
	untestedMoves [][2]int
	testedMoves   [][2]int
	parent        *Node
	children      []*Node
	bestChild     *Node
}

func NewNode(state *game.Board, parent *Node, moveToGetHere [2]int, stone int) *Node {
	moves := state.GetAllLegalMoves()
	children := make([]*Node, 0, len(moves))
	return &Node{
		id:            moveToGetHere,
		nbVisit:       0,
		nbWins:        0,
		playing:       stone,
		state:         *state,
		untestedMoves: moves,
		testedMoves:   [][2]int{},
		parent:        parent,
		children:      children,
		bestChild:     nil,
	}
}

func (n *Node) Simulate() {
	simulationBoard := game.NewBoard()
	*simulationBoard = n.state.GetValue()
	moves := n.untestedMoves
	whosTurn := n.playing
	for len(moves) > 0 {
		randomIndex := rand.Intn(len(moves))
		win, err := simulationBoard.Play(moves[randomIndex][0], moves[randomIndex][1], whosTurn)
		if err != nil {
			n.Backpropagate(1)
			return
		}
		if win != 0 {
			n.Backpropagate(win)
			return
		}
		whosTurn *= -1
		moves[randomIndex] = moves[len(moves)-1]
		moves = moves[:len(moves)]
	}
	n.Backpropagate(1)
}

func (n *Node) Backpropagate(winner int) {
	for n.parent != nil {
		n = n.parent
		n.nbVisit++
		n.nbWins -= winner
		// TODO
	}
}

func NewTree(initialState *game.Board, whosTurn int) *Tree {
	return &Tree{NewNode(initialState, nil, [2]int{}, whosTurn)}
}

func (t *Tree) GetBestMove() [2]int {
	randomIndex := rand.Intn(len(t.root.testedMoves))
	return t.root.testedMoves[randomIndex]
}

func (t *Tree) GetLeaf() *Node {
	node := t.root
	for len(node.untestedMoves) == 0 {
		node = node.bestChild
	}
	randomIndex := rand.Intn(len(node.untestedMoves))
	newBoard := game.NewBoard()
	*newBoard = node.state.GetValue()
	newNode := NewNode(newBoard, node, node.untestedMoves[randomIndex], node.playing*-1)
	_, _ = newNode.state.Play(newNode.id[0], newNode.id[1], newNode.playing)
	node.testedMoves = append(node.testedMoves, node.untestedMoves[randomIndex])
	node.untestedMoves[randomIndex] = node.untestedMoves[len(node.untestedMoves)-1]
	node.untestedMoves = node.untestedMoves[:len(node.untestedMoves)]
	node.children = append(node.children, newNode)
	return newNode
}

func (t *Tree) ComputeMCTS(c chan [2]int, stopChan chan struct{}, playChan chan struct{}) {
	rand.Seed(time.Now().UnixNano())
	for {
		select {
		case <-stopChan:
			return
		case <-playChan:
			bestMove := t.GetBestMove()
			//t.root.parent = nil
			c <- bestMove
		case playedMoves := <-c:
			for i, child := range t.root.children {
				if child.id == playedMoves {
					t.root = t.root.children[i]
					t.root.parent = nil
					break
				}
			}
			newBoard := game.NewBoard()
			*newBoard = t.root.state.GetValue()
			_, _ = newBoard.Play(playedMoves[0], playedMoves[1], t.root.playing*-1)
			t.root = NewNode(newBoard, nil, [2]int{}, t.root.playing*-1)
		default:
			leaf := t.GetLeaf()
			leaf.Simulate()
		}
	}
}
