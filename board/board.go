package board

import "fmt"
import "math"

// Board general board interface that supports MCTS (UCT) algorithm
type Board interface {
	MakeMove(move int)
	TakeMove()
	GetMoves() []int // this should be some other type
	GetResult(playerJM Player) Result // this should also be some other type
	GetPlayerJustMoved() Player
	GetEnemy(playerJM Player) Player
	String() string
}

// TicTacToe implementation of Board interface for the game of Tic Tac Toe
type TicTacToe struct {
	pos [BoardSize]Player
	PlayerJustMoved Player
	history []int
	resultLines []resultLines
}

func (b TicTacToe) String() string {
	// Convert playerJustMoved to playerToMove
	var playerToMove string
	switch b.PlayerJustMoved {  // todo make this a map
	case PlayerO: playerToMove = "X"
	case PlayerX: playerToMove = "O"
	default: playerToMove = "-"
	}

	result := fmt.Sprintf("\n\nPlayer to move %s\n\n", playerToMove)

	for _, rowLine := range b.resultLines[1] {
		var line string
		for _, idx := range rowLine {
			mark := PlayerToString[b.pos[idx]]
			line += fmt.Sprintf("| %s ", mark)
		}
		result += fmt.Sprintf("\t%s|\n", line)
	}
	return result
}

// MakeMove Makes a move on the board and updates related variables (i.e. playerJustMoved etc.)
func (b *TicTacToe) MakeMove(move int) {
	b.PlayerJustMoved = -b.PlayerJustMoved
	b.pos[move] = b.PlayerJustMoved
	b.history = append(b.history, move)
}

// TakeMove Takes back the last made move and updates related variables (i.e. playerJustMoved etc.)
func (b *TicTacToe) TakeMove() {
	if len(b.history) == 0 {
		panic("History is empty")
	}

	lastElementIndex := len(b.history) - 1
	lastElement := b.history[lastElementIndex]
	b.pos[lastElement] = NoPlayer
	b.PlayerJustMoved = -b.PlayerJustMoved

	// remove last element from history -> this only takes a slice
	b.history = b.history[:lastElementIndex]
}

// GetMoves returns all available moves for the current position
func (b *TicTacToe) GetMoves() []int {
	moves := make([]int, 0, BoardSize)

	if b.GetResult(b.PlayerJustMoved) != NoWinner {
		return moves // empty slice of moves
	}

	for i, value := range b.pos {
		if value == NoPlayer {
			moves = append(moves, i)
		}
	}
	return moves
}

func (b *TicTacToe) evaluateLines(lines resultLines, playerJM Player) Result {
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

// GetResult returns the game result from view point of playerJustMoved
func (b *TicTacToe) GetResult(playerJM Player) Result {
	for _, line := range b.resultLines {
		if winner := b.evaluateLines(line, playerJM); winner != NoWinner {
			return winner
		}
	}

	for _, value := range b.pos {
		if value == NoPlayer {
			return NoWinner // there are still available moves
		}
	}

	return Draw  // there is no winner but also no available moves -> Draw
}

// GetPlayerJustMoved returns value of player just moved for Tic Tac Toe
func (b *TicTacToe) GetPlayerJustMoved() Player {
	return b.PlayerJustMoved
}

// GetEnemy Returns the opposite player of the player provided to the func
func (b *TicTacToe) GetEnemy(playerJM Player) Player {
	switch playerJM {
	case PlayerO: return PlayerX
	case PlayerX: return PlayerO
	default: return NoPlayer
	}
}

// CreateNewBoard returns a new instance of a board with default values
func CreateNewBoard() TicTacToe {
	return TicTacToe { PlayerJustMoved: PlayerO, resultLines: getResultLines() }
}