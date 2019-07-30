package uct

import "fmt"
import board "local/gotoe/board"
import "math"

// Node structure to hold information about each node in the MCTS
type Node struct {
	move int
	parent *Node
	childNodes []*Node
	wins float64
	visits float64
	untriedMoves []int
	playerJustMoved board.Player
}

// Update result of game to this node (backpropagate)
func (n *Node) Update(gameResult float64) {
	n.visits += 1.0
	n.wins += gameResult
}

// AddChild adds child node from a given untried move under this node
// note state here can be a pointer to struct
func (n *Node) AddChild(move int, state board.Board) *Node {
	node := Node {
		move: move,
		parent: n,
		untriedMoves: state.GetMoves(),
		playerJustMoved: state.GetPlayerJustMoved(),
	}

	n.childNodes = append(n.childNodes, &node)

	moveIdx := -1
	for idx, m := range n.untriedMoves {
		if m == move {
			moveIdx = idx
		}
	}

	if moveIdx == -1 {
		panic(fmt.Sprintf("Couldn't find move in untried moves %d", move))
	}

	// delete move from untried moves
	// Most efficient way to do this is to swap
	// Unwanted element with the last element
	// and then save slice from start to (excluding) last element
	lastElementIdx := len(n.untriedMoves) - 1
	n.untriedMoves[moveIdx] = n.untriedMoves[lastElementIdx]
	n.untriedMoves = n.untriedMoves[:lastElementIdx]

	return &node
}

func (n *Node) ucb1(node *Node) float64 {
	return (node.wins / node.visits) +
		math.Sqrt(2*math.Log(n.visits) / node.visits)
}

// SelectChild Evaluate all of the node's children using UCB1 formula
// and return most promising one
func (n *Node) SelectChild() *Node {
	bestChildIdx := 0
	bestChildScore := 0.0

	for idx, child := range n.childNodes {
		childScore := n.ucb1(child)
		if childScore > bestChildScore {
			bestChildScore = childScore
			bestChildIdx = idx
		}
	}

	return n.childNodes[bestChildIdx]
}

// CreateRootNode creates a root node for a given board state
func CreateRootNode(state board.Board) Node {
	return Node {
		move: -1, // this is set to an invalid move
		parent: nil,
		untriedMoves: state.GetMoves(),
		playerJustMoved: state.GetPlayerJustMoved(),
	}
}