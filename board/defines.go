package board

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

// PlayerToString maps from Player type to its string representation
var PlayerToString = map[Player]string{
	PlayerO: "O",
	PlayerX: "X",
	NoPlayer: "-",
}

// Result type indicates a predefined game result value used for MCTS (UCT) algorithm
type Result float64
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