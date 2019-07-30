package main

import "fmt"
import board "local/gotoe/board"

func main() {
	b := board.CreateNewBoard()
	fmt.Println(b)
	b.MakeMove(2)
	b.MakeMove(3)
	b.MakeMove(4)
	b.MakeMove(1)
	b.MakeMove(6)
	fmt.Println(b)
	fmt.Println(b.GetResult(b.PlayerJustMoved))
}