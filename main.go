package main

// import "os"
import "fmt"
import "time"
// import "bufio"
// import "strconv"

import board "local/gotoe/board"
import uct "local/gotoe/uct"

func main() {
	b := board.CreateNewBoard()
	start := time.Now()
	fmt.Println(uct.GetEngineMoveFast(&b, 100000))
	fmt.Println(time.Since(start))

	// scanner := bufio.NewScanner(os.Stdin)

	// for b.GetResult(b.GetPlayerJustMoved()) == board.NoWinner {
	// 	if b.PlayerJustMoved == board.PlayerO {
	// 		availableMoves := b.GetMoves()
	// 		fmt.Printf("Enter move (available %v): ", availableMoves)
	// 		scanner.Scan()
	// 		m := scanner.Text()
	// 		move, err := strconv.Atoi(m)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		if move < 0 || move >= board.BoardSize {
	// 			panic(
	// 				fmt.Sprintf("Invalid move %d -> available moves: %v", move, availableMoves))
	// 		}
	// 		b.MakeMove(move)
	// 		fmt.Println(b)
	// 	} else {
	// 		fmt.Println("Engine thinking...")
	// 		move := uct.GetEngineMove(&b, 10000)
	// 		b.MakeMove(move)
	// 		fmt.Printf("Engine makes move %d\n", move)
	// 		fmt.Println(b)
	// 	}
	// }
	// switch b.GetResult(board.PlayerO) {
	// case board.Win:
	// 	fmt.Println("O wins")
	// case board.Draw:
	// 	fmt.Println("Draw!")
	// case board.Loss:
	// 	fmt.Println("X wins")
	// }
}