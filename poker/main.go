package main

import (
	"fmt"
)

func main() {
	fmt.Println("ok")

	cards := NewCards()
	cards.Shuffle()
	cards.PrintCards()
}
