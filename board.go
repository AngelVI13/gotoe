package main

import "math"

// Player type to indicate a player -> X, O, NoPlayer
type Player int8
const (
	PlayerO Player = -1
	NoPlayer Player = 0
	PlayerX Player = 1
)

type Result float32
const (
	NoWinner Result = -1.0
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
	makeMove(move int)
	takeMove()
	getMoves() []int // this should be some other type
	getResult(playerJM Player) Result // this should also be some other type
	String() string // todo
}

type TicTacToe struct {
	pos [BoardSize]Player
	playerJustMoved Player
	history []int
	resultLines []ResultLines
}

func (b *TicTacToe) makeMove(move int) {
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

	// remove last element from history -> this only takes a slice
	b.history = b.history[:lastElementIndex]
}

func (b *TicTacToe) getMoves() []int {
	// todo add check for result

	moves := make([]int, 0, BoardSize)
	for i, value := range b.pos {
		if value == NoPlayer {
			moves = append(moves, i)
		}
	}
	return moves
}

func (b *TicTacToe) evaluateLines(lines ResultLines, playerJM Player) Result {
	for _, line := range lines {
		if line[0] < 0 {
			// if an element from the result line is negative
			// means we are checking diagonal lists which have
			// values only for the first two elements
			continue
		}

		var sum int8
		for _, idx := range line {
			sum += int8(b.pos[idx])
		}
		// the first element of the line would also be the winner

		// todo slow conversion ?
		if int(math.Abs(float64(sum))) == Rows {
			potentialWinner := b.pos[line[0]]

			result := Loss
			if potentialWinner == playerJM {
				result = Win
			}
			return result
		}
	}
	return NoWinner
}

func (b *TicTacToe) getResult(playerJM Player) Result {
	for _, line := range b.resultLines {
		if winner := b.evaluateLines(line, playerJM); winner != NoWinner {
			return winner
		}
	}

	for _, value := range b.pos {
		if value != NoPlayer {
			return NoWinner // there are still available moves
		}
	}

	return Draw  // there is no winner but also no available moves -> Draw
}

type ResultLines [][Rows-1]int

func getColumnArray() ResultLines {
	result := make(ResultLines, Rows-1)
	for i := 0; i < Rows; i++ {
		idx := 0
		for j := i; j < BoardSize; j += i {
			result[i][idx] = j
			idx++
		}
	}
	return result
}

func getRowArray(columnArr ResultLines) ResultLines {
	result := make(ResultLines, Rows-1)

	for i := 0; i < Rows; i++ {
		for j := 0; j < Rows; j++ {
			result[i][j] = columnArr[j][i]
		}
	}
	return result
}

func getDiagonalArray(rowArray ResultLines) ResultLines {
	result := make(ResultLines, Rows-1)

	for i := 0; i < Rows; i++ {
		for j := 0; j < Rows; j++ {
			// Only 2 diagonals will ever be, irrespective of number of rows
			// Switch is needed to specify for the rest of the
			// elements -> set to invalid values so as not to be
			// processed as normal result lines
			switch j {
			case 0:
				result[j][i] = rowArray[i][i]  // left diagonal -> 0,0| 1,1| 2,2
			case 1:
				result[j][i] = rowArray[i][Rows-i-1]  // right diagonal 0,2 | 1,1| 2,0
			default:
				result[j][i] = -1
				// Additional -1 is needed to convert from size to index
			}
		}
	}
	return result
}

func getResultLines() []ResultLines {
	result := make([]ResultLines, Rows-1)
	var columns = getColumnArray()
	var rows = getRowArray(columns)
	var diags = getDiagonalArray(rows)

	result[0] = columns
	result[1] = rows
	result[2] = diags
	return result
}