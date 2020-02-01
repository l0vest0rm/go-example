package main

import "fmt"

func batchTrain(batch int) {
	players := []int{1, 1, 1, 1, 1}
	totalScores := make([]int, len(players), len(players))

	for i := 0; i < batch; i++ {
		game := NewRedTen(players)
		scores := game.Run()
		for j := 0; j < len(scores); j++ {
			totalScores[j] += scores[j]
		}
		if i%100 == 0 {
			fmt.Printf("batch%d,totalScores:%v", i, totalScores)
		}
	}
	fmt.Println("final totalScores:", totalScores)
}

func main() {
	batchTrain(1000)
}
