package main

import "fmt"

func main() {
	s := make([]int, 0)
	fmt.Println(s)
	s = append(s, 1, 2, 3, 4, 5)
	fmt.Println(s)
	s = s[:len(s)-1]
	fmt.Println(s)

}