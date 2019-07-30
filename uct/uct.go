package uct

import "fmt"
import "sort"
import "math/rand"
import board "local/gotoe/board"

type rankedMove struct {
	move int
	score float64
}

func uct(rootstate board.Board, itermax int) float64 {
	rootnode := CreateRootNode(rootstate)

	state := rootstate
	for i := 0; i < itermax; i++ {
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
	return bestMove.wins / bestMove.visits
}

// GetEngineMove returns the best move found by the UCT
// Todo return int is not very flexible?
func GetEngineMove(state board.Board, simulations int) int {
	availableMoves := state.GetMoves()
	simPerMove := simulations / len(availableMoves)

	// todo BoardSize + 1 is not flexible
	bestMove := rankedMove{ move: board.BoardSize + 1, score: 1.0 }
	for _, move := range availableMoves {
		b := state  // todo does this copy ? or points
		b.MakeMove(move)

		/* Check for immediate result after this make_move()
		   It is possible the game is already over by this point
		   in which the value of the move should be immediately computed and
		   put in the result from the view point of the enemy
		   since here moves are evaluated from that viewpoint */
		var score float64
		enemy := b.GetEnemy(b.GetPlayerJustMoved())
		gameResult := b.GetResult(enemy)

		if gameResult != board.NoWinner {
			score = float64(gameResult) / 1.0
		} else {
			score = uct(b, simPerMove)
		}

		fmt.Printf("Move: %d: %f\n", move, score)
		// here the move_score refers to the best enemy reply
    	// therefore we want to minimize that i.e. chose the move
    	// which leads to the lowest scored best enemy reply
		if score < bestMove.score {
			bestMove.score = score
			bestMove.move = move
		}
	}
	return bestMove.move
}