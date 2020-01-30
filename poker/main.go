package main

import (
	"fmt"
)

func main() {
	fmt.Println("ok")
	game := NewGame(RED_TEN, 5)
	game.Cards.PrintCards()
}
