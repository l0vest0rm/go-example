package main

import (
	"fmt"
	"net/http"
)

func servStatic(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
}

func runRounds(players []int, rounds int) []int {
	totalScores := make([]int, len(players), len(players))

	for i := 0; i < rounds; i++ {
		game := NewRedTen(players)
		scores := game.CmdRun()
		for j := 0; j < len(scores); j++ {
			totalScores[j] += scores[j]
		}
		if i%2000 == 0 {
			fmt.Printf("\nround%d,totalScores:%v", i, totalScores)
		}
	}
	return totalScores
}

func batchTrain(batch int, rounds int) {
	players := []int{3, 2, 2, 2, 2}
	finalScores := make([][]int, batch, batch)

	for i := 0; i < batch; i++ {
		fmt.Printf("\nbatch %d begin", i)
		finalScores[i] = runRounds(players, rounds)
	}

	fmt.Println("\nfinal totalScores:")
	for i := 0; i < batch; i++ {
		fmt.Println(finalScores[i])
	}

}

func humanCmdPlay() {
	players := []int{1, 2, 2, 2, 2}

	game := NewRedTen(players)
	scores := game.CmdRun()
	fmt.Printf("scores:%v", scores)
}

func main() {
	webPlay()
	//humanPlay()
	//batchTrain(10, 10000)
}
