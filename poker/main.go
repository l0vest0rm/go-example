package main

import (
	"fmt"
	"net/http"
)

func servStatic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request")
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
}

func batchTrain(batch int) {
	players := []int{1, 1, 1, 1, 1}
	totalScores := make([]int, len(players), len(players))

	for i := 0; i < batch; i++ {
		game := NewRedTen(players)
		scores := game.Run()
		for j := 0; j < len(scores); j++ {
			totalScores[j] += scores[j]
		}
		if i%1000 == 0 {
			fmt.Printf("batch%d,totalScores:%v", i, totalScores)
		}
	}
	fmt.Println("final totalScores:", totalScores)
}

func humanCmdPlay() {
	players := []int{0, 1, 1, 1, 1}

	game := NewRedTen(players)
	scores := game.Run()
	fmt.Printf("scores:%v", scores)
}

func main() {
	webPlay()
	//humanPlay()
	//batchTrain(100000)
}
