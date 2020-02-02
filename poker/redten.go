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
	CmdRun() []int
	CalcScores() []int
	Vals() []int
	RemainCards(playerId int) []int
	NextPlayer(playerId int) int
	PlayerHand(playerId int, candidates []int) []int
	PrintPlayersRemainCards()
}

// IStrategy 策略接口
type IStrategy interface {
	Hand(playerId int, init bool, hands []*Hand, remainCards []int, candidates []int) []int
}

type RedTen struct {
	vals           []int
	playerNum      int
	totalRedTenNum int
	remainCards    CardsVals //整体剩余的卡牌(由大到小排序)
	hands          []*Hand
	players        []*Player
	preHandPlayer  int  //上一个出过牌的人(没出牌不算)
	checkDuty      bool //是否检查蹲
	dutyPlayer     int  //要求蹲的人
	rank           int  //第一个走的
}

type Player struct {
	initCards   CardsVals
	remainCards CardsVals
	strategy    IStrategy
	redTenCnt   int
	rank        int
}

type Hand struct {
	init     bool
	playerId int
	cards    []int
}

type HumanWebPlay struct {
}

type HumanCmdPlay struct {
}

type Robot1 struct {
	redTen *RedTen
}

type Robot2 struct {
	redTen *RedTen
}

// 出牌策略
const (
	GEN_PAI     = 1 //跟拍
	KANG_ZHU    = 2 //抗住
	JINLIANG_DA = 4 //尽量大
	ZUI_DA      = 5 //最大
)

var (
	//值和实际大小映射
	valsMap = []int{
		14, 15, 3, 4, 5, 6, 7, 8, 9, 20, 11, 12, 13,
		14, 15, 3, 4, 5, 6, 7, 8, 9, 20, 11, 12, 13,
		14, 15, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		14, 15, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		30, 25}
	strategys = []func(*RedTen) IStrategy{
		NewHumanWebPlay,
		NewHumanCmdPlay,
		NewRobot1,
		NewRobot2,
	}
)

func NewRedTen(players []int) IGame {
	vals := NewCards()
	t := &RedTen{playerNum: len(players),
		vals:          vals,
		hands:         make([]*Hand, 0),
		players:       make([]*Player, 0),
		preHandPlayer: -1,
		checkDuty:     false,
		dutyPlayer:    -1,
		rank:          1,
	}
	t.ModVals()
	Shuffle(t.vals)
	if t.playerNum == 4 || t.playerNum == 5 {
		t.totalRedTenNum = 2
	}

	for i := 0; i < t.playerNum; i++ {
		t.players = append(t.players, &Player{initCards: make([]int, 0), strategy: strategys[players[i]](t)})
	}

	//发牌
	t.Dispacther()

	return t
}

func (t *RedTen) Vals() []int {
	return t.vals
}

func (t *RedTen) RemainCards(playerId int) []int {
	return t.players[playerId].remainCards
}

// ModVals 修改值
func (t *RedTen) ModVals() {
	l := len(t.vals)
	for i := 0; i < l; i++ {
		if t.vals[i] == 2 || t.vals[i] == 15 || t.vals[i] == 28 || t.vals[i] == 41 {
			if t.playerNum == 5 {
				t.vals[i] = t.vals[l-1]
				l--
			}
		}

	}

	t.vals = t.vals[:l]

	t.remainCards = make([]int, len(t.vals))
	copy(t.remainCards, t.vals)
	sort.Sort(sort.Reverse(t.remainCards))
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
		t.players[i].redTenCnt = redTenCnt(t.players[i].initCards)
		//fmt.Println(i, t.players[i])
	}
}

func isRedTen(card int) bool {
	if card == 9 || card == 22 {
		return true
	} else {
		return false
	}
}

func redTenCnt(cards []int) int {
	redTenCnt := 0
	for _, card := range cards {
		if isRedTen(card) {
			redTenCnt += 1
		}
	}

	return redTenCnt
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
			Printf("\nplayer%d,有%d张牌:", i, len(t.players[i].remainCards))
			PrintCards(t.players[i].remainCards)
		} else {
			Printf("\nplayer%d,有%d张牌", i, n)
		}
	}
}

func (t *RedTen) PrintPlayersInitCards() {
	for i := 0; i < t.playerNum; i++ {
		Printf("\nplayer%d,rank:%d,", i, t.players[i].rank)
		PrintCards(t.players[i].initCards)
	}
}

func (t *RedTen) CalcScores() []int {
	scores := make([]int, t.playerNum, t.playerNum)

	hasRedTenScore := 0
	noRedTenScore := 0
	loserCnt := 0
	totalLoserRedTenCnt := 0

	//找大旗
	first := t.findRank(1)
	firstRedTenCnt := t.players[first].redTenCnt
	//看抓到几个
	for n := t.playerNum; n > 0; n-- {
		loser := t.findRank(n)
		if t.players[loser].redTenCnt == firstRedTenCnt {
			//同伙
			break
		}
		loserCnt += 1
		totalLoserRedTenCnt += t.players[loser].redTenCnt
	}

	switch loserCnt {
	case 0:
	case 1:
		if firstRedTenCnt > 0 {
			//红十先走
			hasRedTenScore = 3
			noRedTenScore = -2
		} else {
			//红十被抓
			if totalLoserRedTenCnt == 1 {
				hasRedTenScore = -3
				noRedTenScore = 2
			} else {
				//双红十被抓
				hasRedTenScore = -32
				noRedTenScore = 8
			}
		}
	case 2:
		if firstRedTenCnt > 0 {
			//红十先走
			hasRedTenScore = 6
			noRedTenScore = -4
		} else {
			//红十被抓
			hasRedTenScore = -9
			noRedTenScore = 6
		}
	case 3:
		//红十先走
		hasRedTenScore = 9
		noRedTenScore = -6
	case 4:
		//红十先走
		hasRedTenScore = 32
		noRedTenScore = -8
	default:
		panic(fmt.Sprintf("loserCnt wrong:%d", loserCnt))
	}

	for i := 0; i < t.playerNum; i++ {
		if t.players[i].redTenCnt > 0 {
			scores[i] = hasRedTenScore
		} else {
			scores[i] = noRedTenScore
		}
	}
	return scores
}

func (t *RedTen) findRank(rank int) int {
	for i := 0; i < t.playerNum; i++ {
		if t.players[i].rank == rank {
			return i
		}
	}

	return -1
}

// Run 开始运行，假设player0是真人，其它是机器人
func (t *RedTen) CmdRun() []int {
	i := -1
	for {
		i = t.NextPlayer(i)
		if i == -1 {
			Println("\ngame over")
			t.PrintPlayersInitCards()
			return t.CalcScores()
		}

		if i == 0 {
			t.PrintPlayersRemainCards()
		}

		t.PlayerHand(i, nil)
	}
}

func (t *RedTen) chooseBeginPlayer() int {
	//找红4在谁那(假设5人)
	for i := 0; i < t.playerNum; i++ {
		for j := 0; j < len(t.players[i].remainCards); j++ {
			if t.players[i].remainCards[j] == 3 {
				return i
			}
		}
	}

	panic("没有找到起始出牌玩家")
}

func (t *RedTen) NextPlayer(playerId int) int {
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
func (t *RedTen) PlayerHand(playerId int, candidates []int) []int {
	init := false
	if t.preHandPlayer == -1 || t.preHandPlayer == playerId {
		init = true
	}

	if t.dutyPlayer == playerId {
		Printf("\nplayer%d 蹲我", playerId)
		init = true
		t.dutyPlayer = -1
	}

	if t.checkDuty {
		t.checkDuty = false
		t.dutyPlayer = playerId
	}

	cards := t.players[playerId].strategy.Hand(playerId, init, t.hands, t.players[playerId].remainCards, candidates)
	if cards == nil {
		return cards
	}

	Printf("player%d hand:", playerId)
	PrintCards(cards)

	t.RecordHand(playerId, cards)

	t.preHandPlayer = playerId
	t.dutyPlayer = -1
	//看是否出完了
	if len(t.players[playerId].remainCards) == 0 {
		Printf("\nplayer%d 走了", playerId)
		t.players[playerId].rank = t.rank
		t.rank += 1
		t.checkDuty = true
	}

	return cards
}

func (t *RedTen) RecordHand(playerId int, cards []int) {
	var err error
	hand := &Hand{init: true, playerId: playerId, cards: cards}
	t.hands = append(t.hands, hand)
	player := t.players[playerId]
	for _, card := range cards {
		player.remainCards, err = removeCard(player.remainCards, card)
		if err != nil {
			fmt.Println("not found in player cards", playerId, card)
		}

		t.remainCards, err = removeCard(t.remainCards, card)
		if err != nil {
			fmt.Println("not found in total remain cards", playerId, card)
		}
	}
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

func NewHumanWebPlay(redTen *RedTen) IStrategy {
	player := &HumanWebPlay{}
	return player
}

// 人出牌
func (t *HumanWebPlay) Hand(playerId int, init bool, hands []*Hand, remainCards []int, candidates []int) []int {
	return candidates
}

func NewHumanCmdPlay(redTen *RedTen) IStrategy {
	player := &HumanCmdPlay{}
	return player
}

// 人出牌
func (t *HumanCmdPlay) Hand(playerId int, init bool, hands []*Hand, remainCards []int, candidates []int) []int {
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
		return t.Hand(playerId, init, hands, remainCards, candidates)
	}

	return vals
}

func NewRobot1(redTen *RedTen) IStrategy {
	player := &Robot1{redTen: redTen}
	return player
}

// 机器人出牌
func (t *Robot1) Hand(playerId int, init bool, hands []*Hand, remainCards []int, candidates []int) []int {
	var cards []int
	typs := arrangeRemainCards(remainCards)
	PrintArrangeRemainCards(typs)
	if init {
		cards = findMultiSame(remainCards, remainCards[0])
		return cards
	}

	preHand := hands[len(hands)-1]
	/*if isCollaborator(t.redTen, playerId, preHand.playerId) && valsMap[preHand.cards[0]] > 9 {
		return nil
	}*/

	//如果是上家出牌，比较大时不要了
	if isPrePlayer(playerId, preHand.playerId, t.redTen.playerNum) && valsMap[preHand.cards[0]] > 12 {
		return nil
	}

	//优先不拆
	cards = findJustBiggerN2(typs[len(preHand.cards)-1], preHand.cards[0])
	if cards == nil {
		cards = unpackMore(typs, preHand.cards[0], len((preHand.cards)))
	}

	return cards
}

func NewRobot2(redTen *RedTen) IStrategy {
	player := &Robot2{redTen: redTen}
	return player
}

// 机器人出牌
func (t *Robot2) Hand(playerId int, init bool, hands []*Hand, remainCards []int, candidates []int) []int {
	var cards []int
	typs := arrangeRemainCards(remainCards)
	PrintArrangeRemainCards(typs)
	if init {
		cards = findMultiSame(remainCards, remainCards[0])
		return cards
	}

	preHand := hands[len(hands)-1]
	//如果已确认是同伙，比较大时不要了
	/*if isCollaborator(t.redTen, playerId, preHand.playerId) && valsMap[preHand.cards[0]] > 9 {
		return nil
	}*/

	//如果是上家出牌,上家大于3张时，比较大时不要了
	if isPrePlayer(playerId, preHand.playerId, t.redTen.playerNum) && valsMap[preHand.cards[0]] > 12 {
		return nil
	}

	//优先不拆
	cards = findJustBiggerN2(typs[len(preHand.cards)-1], preHand.cards[0])
	if cards == nil {
		cards = unpackMore(typs, preHand.cards[0], len((preHand.cards)))
	}

	return cards
}

//整理手牌
func arrangeRemainCards(remainCards []int) [][][]int {
	typs := make([][][]int, 4, 4)
	n := 0
	typs[0] = make([][]int, 0, 0)
	typs[0] = append(typs[0], []int{remainCards[0]})
	for i := 1; i < len(remainCards); i++ {
		if valsMap[remainCards[i-1]] != valsMap[remainCards[i]] {
			n = 0
			typs[n] = append(typs[n], []int{remainCards[i]})
			continue
		}

		//值一样的牌，升级
		l := len(typs[n])
		cards := typs[n][l-1]
		typs[n] = typs[n][:l-1]
		cards = append(cards, remainCards[i])
		n += 1
		if typs[n] == nil {
			typs[n] = make([][]int, 0, 0)
		}
		typs[n] = append(typs[n], cards)

	}

	return typs
}

func PrintArrangeRemainCards(typs [][][]int) {
	for i := 0; i < len(typs); i++ {
		Printf("\n牌型%d:", i)
		for j := 0; j < len(typs[i]); j++ {
			PrintCards(typs[i][j])
			Printf(";")
		}
	}
}

func findJustBiggerN2(cardss [][]int, card int) (vals []int) {
	//找一个刚好大过上家的
	for i := 0; i < len(cardss); i++ {
		if valsMap[card] >= valsMap[cardss[i][0]] {
			continue
		}
		vals = make([]int, len(cardss[i]))
		copy(vals, cardss[i])
		break
	}

	Printf("\nfindJustBiggerN2:")
	PrintCards(vals)

	return vals
}

//拆牌
func unpackMore(typs [][][]int, card int, l int) []int {
	for i := l; i < 4; i++ {
		cards := findJustBiggerN2(typs[i], card)
		if cards != nil {
			return cards[:l]
		}
	}

	return nil
}

func isPrePlayer(playerId int, prePlayerId int, playerNum int) bool {
	if (prePlayerId+1)%playerNum == playerId {
		return true
	}
	return false
}

func isCollaborator(redTenCnt int, hands []*Hand) bool {
	return redTen.players[playerId].redTenCnt == redTen.players[preHandPlayerId].redTenCnt
}
