package main

// Player type to indicate a player -> X, O, NoPlayer
type Player int8
const (
	PlayerO Player = -1
	NoPlayer Player = 0
	PlayerX Player = 1
)

type Result float32
const (
	Loss Result = 0.0
	Draw Result = 0.5
	Win Result = 1.0
)

const (
	Rows = 3
	BoardSize = Rows * Rows
)

// Board general board interface that supports MCTS (UCT) algorithm
type Board interface {
	makeMove(move int8)
	takeMove()
	getMoves() []int8 // this should be some other type
	getResult(playerJM Player) Result // this should also be some other type
	String() string
}

type TicTacToe struct {
	pos [BoardSize]Player
	playerJustMoved Player
	history []int8
}

func (b *TicTacToe) makeMove(move int8) {
	b.playerJustMoved = -b.playerJustMoved
	b.pos[move] = b.playerJustMoved
	b.history = append(b.history, move)
}

func (b *TicTacToe) takeMove() {
	if len(b.history) == 0 {
		panic("History is empty")
	}

	lastElementIndex := len(b.history) - 1
	lastElement := b.history[lastElementIndex]
	b.pos[lastElement] = NoPlayer
	b.playerJustMoved = -b.playerJustMoved

	// remove last element from history
	b.history = b.history[:lastElementIndex]
}
