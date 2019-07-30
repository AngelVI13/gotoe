package main

import "fmt"
import board "local/gotoe/board"
import uct "local/gotoe/uct"

func main() {
	b := board.CreateNewBoard()
	fmt.Println(uct.GetEngineMove(&b, 10000))
}