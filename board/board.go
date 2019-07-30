package board

import "fmt"
import "math"

// Player type to indicate a player -> X, O, NoPlayer
type Player int8
const (
	// PlayerO enum value to indicate player with mark O
	PlayerO Player = -1
	// NoPlayer enum value to indicate a position with no mark on it
	NoPlayer Player = 0
	// PlayerX enum value to indicate player with mark X
	PlayerX Player = 1
)

// Result type indicates a predefined game result value used for MCTS (UCT) algorithm
type Result float32
const (
	// NoWinner value to indicate no winner
	NoWinner Result = -1.0
	// Loss value to indicate a loss for the playerJustMoved
	Loss Result = 0.0
	// Draw value to indicate a drawn game
	Draw Result = 0.5
	// Win value to indicate a win for the playerJustMoved
	Win Result = 1.0
)

const (
	// Rows is the number of rows and columns for the TicTacToe board
	Rows = 3
	// BoardSize is the size of the board field i.e. total number of squares
	BoardSize = Rows * Rows
)

// Board general board interface that supports MCTS (UCT) algorithm
type Board interface {
	MakeMove(move int)
	TakeMove()
	GetMoves() []int // this should be some other type
	GetResult(playerJM Player) Result // this should also be some other type
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
			var mark string
			switch b.pos[idx] { // todo make this a map
			case PlayerO: mark = "O"
			case PlayerX: mark = "X"
			default: mark = "-"
			}

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
		if value != NoPlayer {
			return NoWinner // there are still available moves
		}
	}

	return Draw  // there is no winner but also no available moves -> Draw
}

// CreateNewBoard returns a new instance of a board with default values
func CreateNewBoard() TicTacToe {
	return TicTacToe { PlayerJustMoved: PlayerO, resultLines: getResultLines() }
}

type resultLines [][Rows]int

func getColumnArray() resultLines {
	result := make(resultLines, Rows)

	for i := 0; i < Rows; i++ {
		idx := 0
		for j := i; j < BoardSize; j += Rows {
			result[i][idx] = j
			idx++
		}
	}
	return result
}

func getRowArray(columnArr resultLines) resultLines {
	result := make(resultLines, Rows)

	for i := 0; i < Rows; i++ {
		for j := 0; j < Rows; j++ {
			result[i][j] = columnArr[j][i]
		}
	}
	return result
}

func getDiagonalArray(rowArray resultLines) resultLines {
	result := make(resultLines, Rows)

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

func getResultLines() []resultLines {
	result := make([]resultLines, Rows)
	var columns = getColumnArray()
	var rows = getRowArray(columns)
	var diags = getDiagonalArray(rows)

	result[0] = columns
	result[1] = rows
	result[2] = diags
	return result
}