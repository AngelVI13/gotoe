package main

import "fmt"
import "time"

import board "local/gotoe/board"
import uct "local/gotoe/uct"

func main() {
	b := board.CreateNewBoard()

	start := time.Now()
	fmt.Println(uct.GetEngineMove(&b, 10000))
	fmt.Println(time.Since(start))
}