package main

import "fmt"

// 游戏类型
const (
	RED_TEN = 1
)

// Game 一局游戏
type Game struct {
	playerNum int
	Cards     ICards
	play      IPlay
}

func NewGame(gameType int, playerNum int) *Game {
	var play IPlay

	cards := NewCards()
	if gameType == RED_TEN {
		play = NewRedTen(playerNum)
		cards.SetVals(play.ModVals(cards.Vals()))
	} else {
		panic(fmt.Sprintf("invalid game:%d", gameType))
	}

	cards.Shuffle()
	game := &Game{playerNum: playerNum, Cards: cards, play: play}

	return game
}

// IPlay 游戏局接口
type IPlay interface {
	ModVals(vals []int) []int
	Dispacther()
}
