package main

import (
	"fmt"
)

func test() {
	conf := 0  // 配置、终端默认设置
	bg := 0    // 背景色、终端默认设置
	text := 31 // 前景色、红色
	fmt.Printf("\n %c[%d;%d;%dm%s%c[0m\n\n", 0x1B, conf, bg, text, "testPrintColor", 0x1B)
}

func main() {
	test()
	fmt.Println("ok")
	game := NewRedTen(5)
	PrintCards(game.Vals())
	game.PrintPlayersRemainCards()
	game.Run()

}
