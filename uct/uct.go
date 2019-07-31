package uct

import "fmt"
import "sort"
import "math/rand"
import board "local/gotoe/board"
import "github.com/jinzhu/copier"

type rankedMove struct {
	move int
	score float64
}

type moveScore struct {
	move int
	wins float64
	visits float64
}

func uct(rootstate board.Board, originMove int, itermax int) moveScore {
	rootstate.MakeMove(originMove)
	/* Check for immediate result
		It is possible the game is already over by this point
		in which the value of the move should be immediately computed and
		put in the result from the view point of the enemy
		since here moves are evaluated from that viewpoint */
	result, mScore := isImmediateResult(rootstate, originMove)
	if result == true {
		return mScore
	}

	rootnode := CreateRootNode(rootstate)

	state := rootstate
	for i := 0; i < itermax; i++ {
		// fmt.Println(state)
		node := &rootnode
		movesToRoot := 0

		// Select stage
    	// node is fully expanded and non-terminal
		for len(node.untriedMoves) == 0 && len(node.childNodes) > 0 {
			node = node.SelectChild()
			state.MakeMove(node.move)
			movesToRoot++
		}

		// Expand
    	// if we can expand (i.e. state/node is non-terminal)
		if len(node.untriedMoves) > 0 {
			move := node.untriedMoves[rand.Intn(len(node.untriedMoves))]
			state.MakeMove(move)
			movesToRoot++
			// add child and descend tree
			node = node.AddChild(move, state)
		}

		// Rollout
		//  - this can often be made orders of magnitude quicker
		//    using a state.GetRandomMove() function

		// while state is non-terminal
		for state.GetResult(state.GetPlayerJustMoved()) == board.NoWinner {
			moves := state.GetMoves()
			m := moves[rand.Intn(len(moves))]
			state.MakeMove(m)
			movesToRoot++
		}

		// Backpropagate
    	// backpropagate from the expanded node and work back to the root node
		for node != nil {
			// state is terminal.
      		// Update node with result from POV of node.playerJustMoved
			gameResult := state.GetResult(node.playerJustMoved)
			node.Update(float64(gameResult))
			node = node.parent
		}

		for j := 0; j < movesToRoot; j++ {
			state.TakeMove()
		}
	}

	sort.Slice(rootnode.childNodes, func(i, j int) bool {
		return rootnode.childNodes[i].visits > rootnode.childNodes[j].visits
	})
	// above we sort by descending order -> move with most visits is the first element
	bestMove := rootnode.childNodes[0]
	return moveScore { move: originMove, wins: bestMove.wins, visits: bestMove.visits }
}

type uctArg struct {
	state board.Board
	move int
	simulations int
}

func worker(jobs <-chan uctArg, results chan<- moveScore) {
	for uctArguments := range jobs {
		results <- uct(uctArguments.state, uctArguments.move, uctArguments.simulations)
	}
}

// GetEngineMoveFast returns the best move found by the UCT (computed in parallel)
func GetEngineMoveFast(state board.Board, simulations int) int {
	availableMoves := state.GetMoves()
	numMoves := len(availableMoves)

	if numMoves == 0 {
		panic("Game is already over, can't get engine move for a finished game!")
	}

	simPerMove := simulations / numMoves

	// create channels to share data between goroutines
	jobs := make(chan uctArg, numMoves)
	results := make(chan moveScore, numMoves)

	// spawn workers ready to process data
	for _i := 0; _i < numMoves; _i++ {
		go worker(jobs, results)
	}

	// todo BoardSize + 1 is not flexible
	bestMove := rankedMove{ move: board.BoardSize + 1, score: 1.0 }

	for _, move := range availableMoves {
		// create a copy of the board in order to be sent to the goroutine
		b := board.TicTacToe {}
		copier.Copy(&b, state)

		jobs <- uctArg {
			state: &b,
			move: move,
			simulations: simPerMove,
		}
	}

	close(jobs) // close jobs channel

	var mScore moveScore
	var scoreValue float64
	for _i := 0; _i < numMoves; _i++ {
		mScore = <-results
		scoreValue = mScore.wins/mScore.visits

		fmt.Printf("Move: %d: %.3f -> %f / %f\n", mScore.move, scoreValue, mScore.wins, mScore.visits)
		// here the move_score refers to the best enemy reply
    	// therefore we want to minimize that i.e. chose the move
    	// which leads to the lowest scored best enemy reply
		if scoreValue < bestMove.score {
			bestMove.score = scoreValue
			bestMove.move = mScore.move
		}
	}

	return bestMove.move
}

func isImmediateResult(state board.Board, move int) (result bool, score moveScore) {
	enemy := state.GetEnemy(state.GetPlayerJustMoved())
	gameResult := state.GetResult(enemy)

	if gameResult != board.NoWinner {
		score = moveScore { move: move, wins: float64(gameResult), visits: 1.0 }
		result = true
	}

	return
}

// GetEngineMove returns the best move found by the UCT
func GetEngineMove(state board.Board, simulations int) int {
	availableMoves := state.GetMoves()

	if len(availableMoves) == 0 {
		panic("Game is already over, can't get engine move for a finished game!")
	}

	simPerMove := simulations / len(availableMoves)

	// todo BoardSize + 1 is not flexible
	bestMove := rankedMove{ move: board.BoardSize + 1, score: 1.0 }
	for _, move := range availableMoves {
		b := state  // todo does this copy ? or points

		mScore := uct(b, move, simPerMove)
		scoreValue := mScore.wins/mScore.visits

		fmt.Printf("Move: %d: %.3f -> %f / %f\n", move, scoreValue, mScore.wins, mScore.visits)
		// here the move_score refers to the best enemy reply
    	// therefore we want to minimize that i.e. chose the move
    	// which leads to the lowest scored best enemy reply
		if scoreValue < bestMove.score {
			bestMove.score = scoreValue
			bestMove.move = move
		}
		// take move since b here is not a copy but points to the same position as state
		b.TakeMove()
	}
	return bestMove.move
}