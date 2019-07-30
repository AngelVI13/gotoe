package main

import "fmt"
import board "local/gotoe/board"
import uct "local/gotoe/uct"

func main() {
	b := board.CreateNewBoard()
	rootnode := uct.CreateRootNode(&b)
	fmt.Println(rootnode)
}