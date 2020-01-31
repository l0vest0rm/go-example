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

// IStrategy 策略接口
type IStrategy interface {
	Hand(playerId int, init bool, hands []*Hand, remainCards []int) []int
}

type RedTen struct {
	vals      []int
	playerNum int
	hands     []*Hand
	players   []*Player
}

type Player struct {
	initCards   CardsVals
	remainCards CardsVals
	strategy    IStrategy
}

type Hand struct {
	init     bool
	playerId int
	cards    []int
}

type Robot1 struct {
}

type Human struct {
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

	var strategy IStrategy
	for i := 0; i < t.playerNum; i++ {
		if i == -1 {
			strategy = NewHuman()
		} else {
			strategy = NewRobot1()
		}
		t.players = append(t.players, &Player{initCards: make([]int, 0), strategy: strategy})
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
		t.players[idx].initCards = append(t.players[idx].initCards, t.vals[i])
		idx = (idx + 1) % t.playerNum
	}

	//排序
	for i := 0; i < t.playerNum; i++ {
		sort.Sort(t.players[i].initCards)
		t.players[i].remainCards = make([]int, len(t.players[i].initCards))
		copy(t.players[i].remainCards, t.players[i].initCards)
	}
}

func (t *RedTen) PrintPlayersRemainCards() {
	n := 99
	for i := 0; i < t.playerNum; i++ {
		if len(t.players[i].remainCards) > 3 {
			n = 99
		} else {
			n = len(t.players[i].remainCards)
		}

		if i == 0 {
			fmt.Printf("\nplayer%d,有%d张牌:", i, len(t.players[i].remainCards))
			PrintCards(t.players[i].remainCards)
		} else {
			fmt.Printf("\nplayer%d,有%d张牌", i, n)
		}
	}
}

func (t *RedTen) PrintPlayersInitCards() {
	for i := 0; i < t.playerNum; i++ {
		fmt.Printf("\nplayer%d:", i)
		PrintCards(t.players[i].initCards)
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
			fmt.Println("\ngame over")
			t.PrintPlayersInitCards()
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
			fmt.Printf("\nplayer%d 蹲我", i)
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
			fmt.Printf("\nplayer%d hand:%v", i, ConvertVals2PrintChars(cards))
			prePlayer = i
			dutyPlayer = -1
			//看是否出完了
			if len(t.players[i].remainCards) == 0 {
				fmt.Printf("\nplayer%d 走了", i)
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
	cards := t.players[playerId].strategy.Hand(playerId, init, t.hands, t.players[playerId].remainCards)
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

func convertStr2Val(input string) ([]int, error) {
	cards := strings.Split(input, ",")
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

// 策略

func NewHuman() IStrategy {
	player := &Human{}
	return player
}

// 人出牌
func (t *Human) Hand(playerId int, init bool, hands []*Hand, remainCards []int) []int {
	var input string
	fmt.Println("\n请出牌:")
	fmt.Scanln(&input)
	//fmt.Printf("player%d hand:%s\n", playerId, input)
	if input == "" {
		//跳过不出
		return nil
	}
	vals, err := convertStr2Val(input)
	if err != nil {
		fmt.Println("出牌有误,", err)
		return t.Hand(playerId, init, hands, remainCards)
	}

	return vals
}

func NewRobot1() IStrategy {
	player := &Robot1{}
	return player
}

// 机器人出牌
func (t *Robot1) Hand(playerId int, init bool, hands []*Hand, remainCards []int) []int {
	var cards []int
	if init {
		cards = findMultiSame(remainCards, remainCards[0])
	} else {
		preHand := hands[len(hands)-1]
		cards = findJustBiggerN(remainCards, preHand.cards[0], len(preHand.cards))
	}

	return cards
}
