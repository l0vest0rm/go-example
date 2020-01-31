package main

import (
	"fmt"
)

func main() {
	fmt.Println("ok")
	game := NewRedTen(5)
	PrintCards(game.Vals())
	game.PrintPlayersRemainCards()
	game.Run()
}
