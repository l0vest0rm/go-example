package main

import (
	"fmt"
	"sort"
	"strings"
)

type CardsVals []int

// IGame 游戏局接口
type IGame interface {
	Dispacther()
	Run()
	Vals() []int
	PrintPlayersRemainCards()
}

type RedTen struct {
	vals      []int
	playerNum int
	hands     []*Hand
	players   []*Player
}

type Player struct {
	remainCards CardsVals
}

type Hand struct {
	init     bool
	playerId int
	cards    []int
}

var (
	//值和实际大小映射
	valsMap = []int{90,
		14, 20, 3, 4, 5, 6, 7, 8, 9, 30, 11, 12, 13,
		14, 20, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		14, 20, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		14, 20, 3, 4, 5, 6, 7, 8, 9, 30, 11, 12, 13,
		80}
)

func NewRedTen(playerNum int) IGame {
	vals := NewCards()
	t := &RedTen{playerNum: playerNum,
		vals:    vals,
		hands:   make([]*Hand, 0),
		players: make([]*Player, 0)}
	t.ModVals()
	Shuffle(t.vals)

	for i := 0; i < t.playerNum; i++ {
		t.players = append(t.players, &Player{remainCards: make([]int, 0)})
	}

	//发牌
	t.Dispacther()

	return t
}

func (t *RedTen) Vals() []int {
	return t.vals
}

// ModVals 修改值
func (t *RedTen) ModVals() {
	l := len(t.vals)
	for i := 0; i < l; i++ {
		if t.vals[i] == 3 || t.vals[i] == 16 || t.vals[i] == 29 || t.vals[i] == 42 {
			if t.playerNum == 5 {
				t.vals[i] = t.vals[l-1]
				l--
			}
		}

	}

	t.vals = t.vals[:l]
}

func (t *RedTen) Dispacther() {
	idx := 0
	for i := 0; i < len(t.vals); i++ {
		t.players[idx].remainCards = append(t.players[idx].remainCards, t.vals[i])
		idx = (idx + 1) % t.playerNum
	}

	//排序
	for i := 0; i < t.playerNum; i++ {
		sort.Sort(t.players[i].remainCards)
	}
}

func (t *RedTen) PrintPlayersRemainCards() {
	for i := 0; i < t.playerNum; i++ {
		fmt.Printf("\nplayer%d:", i)
		PrintCards(t.players[i].remainCards)
	}
}

// Run 开始运行，假设player0是真人，其它是机器人
func (t *RedTen) Run() {
	init := false
	i := -1
	prePlayer := i
	checkDuty := false
	dutyPlayer := -1 //上一个走的人
	for {
		i = t.nextPlayer(i)
		if i == -1 {
			fmt.Println("game over")
			break
		}
		if prePlayer == -1 {
			prePlayer = i
		}

		if prePlayer == i {
			init = true
		} else {
			init = false
		}

		if dutyPlayer == i {
			fmt.Printf("player%d 蹲我\n", i)
			init = true
			dutyPlayer = -1
		}

		if checkDuty {
			checkDuty = false
			dutyPlayer = i
		}

		if i == 0 {
			t.PrintPlayersRemainCards()
		}

		cards := t.playerHand(i, init)
		if cards != nil {
			prePlayer = i
			dutyPlayer = -1
			//看是否出完了
			if len(t.players[i].remainCards) == 0 {
				fmt.Printf("player%d 走了\n", i)
				checkDuty = true
			} else {
				checkDuty = false
			}
		}
	}
}

func (t *RedTen) chooseBeginPlayer() int {
	//找红4在谁那(假设5人)
	for i := 0; i < t.playerNum; i++ {
		for j := 0; j < len(t.players[i].remainCards); j++ {
			if t.players[i].remainCards[j] == 4 {
				return i
			}
		}
	}

	panic("没有找到起始出牌玩家")
}

func (t *RedTen) nextPlayer(playerId int) int {
	if playerId == -1 {
		return t.chooseBeginPlayer()
	}

	for i := 0; i < t.playerNum; i++ {
		playerId = (playerId + 1) % t.playerNum
		if len(t.players[playerId].remainCards) > 0 {
			return playerId
		}
	}

	return -1
}

// 出一手牌,如果返回空表示不出
func (t *RedTen) playerHand(playerId int, init bool) []int {
	var cards []int
	if playerId == 0 {
		cards = t.humanHand(playerId)
	} else {
		cards = t.botHand(playerId, init)
	}

	if cards == nil {
		return cards
	}

	hand := &Hand{init: init, playerId: playerId, cards: cards}
	t.hands = append(t.hands, hand)
	player := t.players[playerId]
	for _, card := range cards {
		found := false
		l := len(player.remainCards)
		for i := 0; i < l; i++ {
			if card == player.remainCards[i] {
				found = true
				player.remainCards = append(player.remainCards[:i], player.remainCards[i+1:]...)
				break
			}
		}
		if found == false {
			fmt.Println("not found", playerId, card)
		}
	}

	return cards
}

// 机器人出牌
func (t *RedTen) botHand(playerId int, init bool) []int {
	var cards []int
	if init {
		cards = findMultiSame(t.players[playerId].remainCards, t.players[playerId].remainCards[0])
	} else {
		preHand := t.hands[len(t.hands)-1]
		player := t.players[playerId]
		cards = findJustBiggerN(player.remainCards, preHand.cards[0], len(preHand.cards))
	}

	fmt.Printf("player%d hand:%v\n", playerId, ConvertVals2PrintChars(cards))
	return cards
}

// 人出牌
func (t *RedTen) humanHand(playerId int) []int {
	var input string
	fmt.Println("\n请出牌:")
	fmt.Scanln(&input)
	fmt.Printf("player%d hand:%s\n", playerId, input)
	if input == "" {
		//跳过不出
		return nil
	}
	vals, err := t.convertStr2Val(input)
	if err != nil {
		fmt.Println("出牌有误,", err)
		return t.humanHand(playerId)
	}

	return vals
}

func (t *RedTen) convertStr2Val(input string) ([]int, error) {
	cards := strings.Split(input, ",")
	fmt.Println(len(cards))
	vals := make([]int, 0)
	for i := 0; i < len(cards); i++ {
		val, err := ConvertStr2Val(cards[i])
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}

	return vals, nil
}

func findMultiSame(cards []int, card int) (vals []int) {
	for i := 0; i < len(cards); i++ {
		if valsMap[cards[i]] == valsMap[card] {
			if vals == nil {
				vals = make([]int, 0)
			}
			vals = append(vals, cards[i])
		} else if valsMap[cards[i]] > valsMap[card] {
			break
		}
	}

	return vals
}

func findJustBiggerOne(cards []int, card int) (vals []int) {
	//找一个刚好大过上家的
	l := len(cards)
	for i := 0; i < l; i++ {
		if valsMap[card] < valsMap[cards[i]] {
			if vals == nil {
				vals = make([]int, 0)
			}
			vals = append(vals, cards[i])
			break
		}
	}
	return vals
}

func findJustBiggerN(cards []int, card int, n int) (vals []int) {
	//找一个刚好大过上家的
	l := len(cards)
	sameNum := 0
	preVal := -1
	for i := 0; i < l; i++ {
		if valsMap[card] >= valsMap[cards[i]] {
			continue
		}

		if sameNum == 0 || valsMap[cards[i]] == preVal {
			sameNum += 1
			preVal = valsMap[cards[i]]
		}

		if sameNum == n {
			vals = make([]int, 0)
			for j := 0; j < n; j++ {
				vals = append(vals, cards[i-j])
			}
			break
		}
	}
	return vals
}

//Len()
func (t CardsVals) Len() int {
	return len(t)
}

//Less(): 由小到大排序
func (t CardsVals) Less(i, j int) bool {
	return valsMap[t[i]] < valsMap[t[j]]
}

//Swap()
func (t CardsVals) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
