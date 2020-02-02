package main

import (
	"fmt"
	"net/http"
)

func servStatic(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
}

func batchTrain(batch int) {
	players := []int{3, 2, 2, 2, 2}
	totalScores := make([]int, len(players), len(players))

	for i := 0; i < batch; i++ {
		game := NewRedTen(players)
		scores := game.CmdRun()
		for j := 0; j < len(scores); j++ {
			totalScores[j] += scores[j]
		}
		if i%1000 == 0 {
			fmt.Printf("\nbatch%d,totalScores:%v", i, totalScores)
		}
	}
	fmt.Println("\nfinal totalScores:", totalScores)
}

func humanCmdPlay() {
	players := []int{1, 2, 2, 2, 2}

	game := NewRedTen(players)
	scores := game.CmdRun()
	fmt.Printf("scores:%v", scores)
}

func main() {
	//webPlay()
	//humanPlay()
	batchTrain(10000)
}
