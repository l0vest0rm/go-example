package doudizhu

var (
	valsMap = []int{
		14, 15, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		14, 15, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		14, 15, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		14, 15, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		30, 25}
)

func aviableCandidates(remainCards []int) [][]int {
	candidates := make([][]int, 0)
	candidates = append(candidates, []int{remainCards[0]})
	for i := 1; i < len(remainCards); i++ {
		if valsMap[remainCards[i-1]] != valsMap[remainCards[i]] {
			candidates = append(candidates, []int{remainCards[i]})
			continue
		}

		//值一样的牌，升级拷贝
		cards := candidates[len(candidates)-1]
		//typs[n] = typs[n][:l-1]
		newCards := make([]int, len(cards))
		copy(newCards, cards)
		newCards = append(newCards, remainCards[i])
		candidates = append(candidates, newCards)
	}

	return candidates
}

func aviableBiggerCandidates(cards []int, preHand []int) [][]int {
	//找一个刚好大过上家的
	sameNum := 0
	preVal := -1
	card := preHand[0]
	n := len(preHand)
	candidates := make([][]int, 0)
	for i := 0; i < len(cards); i++ {
		if valsMap[card] >= valsMap[cards[i]] {
			continue
		}

		if valsMap[cards[i]] != preVal {
			sameNum = 1
			preVal = valsMap[cards[i]]
		} else {
			sameNum++
		}

		if sameNum == n {
			vals := make([]int, 0)
			for j := 0; j < n; j++ {
				vals = append(vals, cards[i-j])
			}
			candidates = append(candidates, vals)
		}
	}
	return candidates
}
