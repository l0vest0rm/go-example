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
	i := t.chooseBeginPlayer()
	prePlayer := i
	for {
		if i == 0 {
			t.PrintPlayersRemainCards()
		}
		if prePlayer == i {
			init = true
		} else {
			init = false
		}
		cards := t.playerHand(i, init)
		if cards != nil {
			prePlayer = i
		}
		i = t.nextPlayer(i)
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
	return (playerId + 1) % t.playerNum
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
		cards = []int{t.players[playerId].remainCards[0]}
	} else {
		preHand := t.hands[len(t.hands)-1]
		player := t.players[playerId]
		//找一个刚好大过上家的
		l := len(player.remainCards)
		for i := 0; i < l; i++ {
			if valsMap[preHand.cards[0]] < valsMap[player.remainCards[i]] {
				cards = append(cards, player.remainCards[i])
				break
			}
		}
	}

	fmt.Printf("player%d hand:%v\n", playerId, ConvertVals2Strs(cards))
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
