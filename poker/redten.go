package main

import (
	"fmt"
	"sort"
	"strings"

	"./card"
	"./gomcts"
	"github.com/cnf/structhash"
	"github.com/spaolacci/murmur3"
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
	isYourTurn(playerId int) bool
	PlayerHand(playerId int, candidates []int) ([]int, bool)
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
	lastTurnPlayer int //上一个出过牌的人(没出牌不算)
	rank           int //第一个走的
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
	NO       = 0
	YES      = 1
	NOT_SURE = 2 //不确定
)

const (
	RED_TEN_VALUE = 20
)

var (
	//值和实际大小映射
	valsMap = []int{
		14, 15, 3, 4, 5, 6, 7, 8, 9, RED_TEN_VALUE, 11, 12, 13,
		14, 15, 3, 4, 5, 6, 7, 8, 9, RED_TEN_VALUE, 11, 12, 13,
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
	vals := card.NewCards()
	t := &RedTen{playerNum: len(players),
		vals:           vals,
		hands:          make([]*Hand, 0),
		players:        make([]*Player, 0),
		lastTurnPlayer: -1,
		rank:           1,
	}
	t.ModVals()
	card.Shuffle(t.vals)
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
		t.players[i].redTenCnt = clcRedTenCnt(t.players[i].initCards)
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

func clcRedTenCnt(cards []int) int {
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
			card.PrintCards(t.players[i].remainCards)
		} else {
			Printf("\nplayer%d,有%d张牌", i, n)
		}
	}
}

func (t *RedTen) PrintPlayersInitCards() {
	for i := 0; i < t.playerNum; i++ {
		Printf("\nplayer%d,rank:%d,", i, t.players[i].rank)
		card.PrintCards(t.players[i].initCards)
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

func (t *RedTen) isYourTurn(playerId int) bool {
	if t.lastTurnPlayer == -1 || t.lastTurnPlayer == playerId {
		return true
	}

	return false
}

// 出一手牌,如果返回空表示不出
func (t *RedTen) PlayerHand(playerId int, candidates []int) (cards []int, valid bool) {
	//临时数据保存
	lastTurnPlayer := t.lastTurnPlayer

	init := false
	if t.lastTurnPlayer == -1 || t.lastTurnPlayer == playerId {
		init = true
	}

	//上一个人走了，考虑是否蹲我
	if t.lastTurnPlayer != -1 && len(t.players[t.lastTurnPlayer].remainCards) == 0 {
		Printf("\nplayer%d 蹲我吗", playerId)
		t.lastTurnPlayer = playerId
	}

	//检查候选是否合法
	if candidates != nil && len(candidates) > 0 {
		if init {
			valid = validateSame(candidates)
		} else {
			valid = validateBigger(t.hands[len(t.hands)-1].cards, candidates)
		}
		if !valid {
			//恢复现场
			t.lastTurnPlayer = lastTurnPlayer
			return nil, false
		}
	}

	cards = t.players[playerId].strategy.Hand(playerId, init, t.hands, t.players[playerId].remainCards, candidates)
	if cards == nil || len(cards) == 0 {
		return nil, true
	}

	Printf("player%d hand:", playerId)
	card.PrintCards(cards)

	t.RecordHand(playerId, cards)
	t.lastTurnPlayer = playerId

	//看是否出完了
	if len(t.players[playerId].remainCards) == 0 {
		Printf("\nplayer%d 走了", playerId)
		t.players[playerId].rank = t.rank
		t.rank += 1
	}

	return cards, true
}

func (t *RedTen) RecordHand(playerId int, cards []int) {
	var err error
	hand := &Hand{init: true, playerId: playerId, cards: cards}
	t.hands = append(t.hands, hand)
	player := t.players[playerId]
	for _, c := range cards {
		player.remainCards, err = card.RemoveCard(player.remainCards, c)
		if err != nil {
			fmt.Println("not found in player cards", playerId, c)
		}

		t.remainCards, err = card.RemoveCard(t.remainCards, c)
		if err != nil {
			fmt.Println("not found in total remain cards", playerId, c)
		}
	}
}

func convertStr2Val(input string) ([]int, error) {
	cards := strings.Split(input, ",")
	vals := make([]int, 0)
	for i := 0; i < len(cards); i++ {
		val, err := card.ConvertStr2Val(cards[i])
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

		if valsMap[cards[i]] != preVal {
			sameNum = 1
			preVal = valsMap[cards[i]]
		} else {
			sameNum++
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
	if init {
		cards = findMultiSame(remainCards, remainCards[0])
		return cards
	}

	preHand := hands[len(hands)-1]
	//如果已确认是同伙，比较大时不要了
	ret := isCollaborator(t.redTen.hands, t.redTen.totalRedTenNum, t.redTen.players[playerId].redTenCnt, playerId)
	if ret == YES && valsMap[preHand.cards[0]] > 9 {
		return nil
	}

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
	if init {
		cards = findMultiSame(remainCards, remainCards[0])
		return cards
	}

	preHand := hands[len(hands)-1]
	//如果已确认是同伙，比较大时不要了
	ret := isCollaborator(t.redTen.hands, t.redTen.totalRedTenNum, t.redTen.players[playerId].redTenCnt, playerId)
	if ret == YES && valsMap[preHand.cards[0]] > 9 {
		return nil
	}

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
			card.PrintCards(typs[i][j])
			Printf(";")
		}
	}
}

func findJustBiggerN2(cardss [][]int, c int) (vals []int) {
	//找一个刚好大过上家的
	for i := 0; i < len(cardss); i++ {
		if valsMap[c] >= valsMap[cardss[i][0]] {
			continue
		}
		vals = make([]int, len(cardss[i]))
		copy(vals, cardss[i])
		break
	}

	Printf("\nfindJustBiggerN2:")
	card.PrintCards(vals)

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

func getRedTenPlayerIds(hands []*Hand, totalRedTenNum int, myRedTenCnt int, myPlayerId int) ([]int, bool) {
	playerIds := make([]int, 0)
	redTenCnt := 0
	n := 0

	if myRedTenCnt > 0 {
		redTenCnt += myRedTenCnt
		playerIds = append(playerIds, myPlayerId)
	}

	l := len(hands)
	for i := 0; i < l; i++ {
		if hands[i].playerId == myPlayerId {
			continue
		}

		n = clcRedTenCnt(hands[i].cards)
		if n > 0 {
			redTenCnt += n
			found := false
			for j := 0; j < len(playerIds); j++ {
				if playerIds[j] == hands[i].playerId {
					found = true
				}
			}

			if !found {
				playerIds = append(playerIds, hands[i].playerId)
			}
		}
	}

	if redTenCnt == totalRedTenNum {
		return playerIds, true
	} else {
		return playerIds, false
	}
}

func isCollaborator(hands []*Hand, totalRedTenNum int, myRedTenCnt int, myPlayerId int) int {
	redTenPlayerIds, isAll := getRedTenPlayerIds(hands, totalRedTenNum, myRedTenCnt, myPlayerId)
	preHnaderHasRedTen := false
	l := len(hands)
	for _, playerId := range redTenPlayerIds {
		if playerId == hands[l-1].playerId {
			preHnaderHasRedTen = true
			break
		}
	}

	if preHnaderHasRedTen {
		if myRedTenCnt > 0 {
			return YES
		} else {
			return NO
		}
	}

	//红十都出完了，上家确定没有红十
	if isAll {
		if myRedTenCnt > 0 {
			return NO
		} else {
			return YES
		}
	}

	//不确定性上家是否有红十
	return NOT_SURE
}

//验证牌型有效（都一样）
func validateSame(cards []int) bool {
	if len(cards) == 1 {
		return true
	}

	for i := 1; i < len(cards); i++ {
		if valsMap[cards[i]] != valsMap[cards[i-1]] {
			return false
		}
	}

	return true
}

func validateBigger(preCards []int, cards []int) bool {
	if len(preCards) > len(cards) {
		return false
	}

	if !validateSame(cards) {
		return false
	}

	if len(preCards) == len(cards) {
		return valsMap[cards[0]] > valsMap[preCards[0]]
	}

	//判断是否红十
	if valsMap[cards[0]] == RED_TEN_VALUE {
		return true
	}

	return false
}

func findWinHand(a []int, b []int, preHand []int) []int {
	var candidates [][]int
	if preHand == nil || len(preHand) == 0 {
		//新出
		candidates = aviableCandidates(a)
	} else {
		candidates = aviableBiggerCandidates(a, preHand)
		candidates = append(candidates, []int{})
	}

	//PrintCandidatesCards(candidates)
	for i := 0; i < len(candidates); i++ {
		a1, err := card.RemoveCards(a, candidates[i])
		if err != nil {
			panic(err)
		}

		if len(a1) == 0 {
			//fmt.Printf("\nlen(a1) == 0:%v", ConvertVals2PrintChars(candidates[i]))
			return candidates[i]
		}

		if innerFindWinSolution(a1, b, candidates[i], false) {
			//fmt.Printf("\ncandidatesA:%v", ConvertVals2PrintChars(candidates[i]))
			return candidates[i]
		}
	}

	return nil
}

//返回是否继续出牌
func innerFindWinSolution(a []int, b []int, preHand []int, firstTurn bool) bool {
	var turn []int
	var candidates [][]int
	if firstTurn {
		turn = a
	} else {
		turn = b
	}

	if preHand == nil || len(preHand) == 0 {
		//新出
		candidates = aviableCandidates(turn)
	} else {
		candidates = aviableBiggerCandidates(turn, preHand)
		candidates = append(candidates, []int{})
	}

	if firstTurn {
		for i := 0; i < len(candidates); i++ {
			a1, err := card.RemoveCards(a, candidates[i])
			if err != nil {
				panic(err)
			}

			if len(a1) == 0 {
				//fmt.Printf("\n innerFindWinSolution,len(a1) == 0:%v", ConvertVals2PrintChars(candidates[i]))
				return true
			}

			if innerFindWinSolution(a1, b, candidates[i], false) {
				//fmt.Printf("\n innerFindWinSolution,candidatesA:%v", ConvertVals2PrintChars(candidates[i]))
				return true
			}
		}

		return false
	}

	//得要所有的b出牌策略都赢才算赢
	for i := 0; i < len(candidates); i++ {
		b1, err := card.RemoveCards(b, candidates[i])
		if err != nil {
			panic(err)
		}

		if len(b1) == 0 {
			return false
		}

		if !innerFindWinSolution(a, b1, candidates[i], true) {
			return false
		}
	}

	//fmt.Printf("\n innerFindWinSolution,anyB,A:%v,B:%v,preHand:%v", ConvertVals2PrintChars(a), ConvertVals2PrintChars(b), ConvertVals2PrintChars(preHand))
	return true
}

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

func PrintCandidatesCards(candidates [][]int) {
	for i := 0; i < len(candidates); i++ {
		card.PrintCards(candidates[i])
		Printf(";")
	}
}

func aviableCandidatesOld(remainCards []int) [][][]int {
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

		//值一样的牌，升级拷贝
		l := len(typs[n])
		cards := typs[n][l-1]
		//typs[n] = typs[n][:l-1]
		newCards := make([]int, len(cards))
		copy(newCards, cards)
		newCards = append(newCards, remainCards[i])
		n += 1
		if typs[n] == nil {
			typs[n] = make([][]int, 0, 0)
		}
		typs[n] = append(typs[n], newCards)
	}

	return typs
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

// mcts算法相关

// RedTenGameAction - action on a tic tac toe board game
type RedTenGameAction struct {
	hand       []int
	nextToMove int8
}

func (a RedTenGameAction) Hash() uint64 {
	return 0
}

func (t *RedTenGameAction) ApplyTo(s gomcts.GameState) gomcts.GameState {
	var err error
	state := s.(*RedTenGameState)
	if state.nextToMove != t.nextToMove {
		panic("wrong turn")
	}

	newState := &RedTenGameState{
		a:          state.a,
		b:          state.b,
		preHand:    t.hand,
		nextToMove: -1 * state.nextToMove,
	}

	if t.nextToMove == 1 {
		newState.a, err = card.RemoveCards(newState.a, t.hand)
	} else {
		newState.b, err = card.RemoveCards(newState.b, t.hand)
	}

	if err != nil {
		panic(err)
	}

	return newState
}

type RedTenGameState struct {
	a          []int
	b          []int
	preHand    []int
	nextToMove int8
}

func (a RedTenGameState) Hash() uint64 {
	return 0
}

func (t *RedTenGameState) EvaluateGame() (gomcts.GameResult, bool) {
	if len(t.a) == 0 {
		return 1, true
	} else if len(t.b) == 0 {
		return -1, true
	} else {
		return 0, false
	}
}

func (t *RedTenGameState) GetLegalActions() []gomcts.Action {
	var turn []int
	var candidates [][]int

	if t.nextToMove == 1 {
		turn = t.a
	} else {
		turn = t.b
	}

	if t.preHand == nil || len(t.preHand) == 0 {
		//新出
		candidates = aviableCandidates(turn)
	} else {
		candidates = aviableBiggerCandidates(turn, t.preHand)
		candidates = append(candidates, []int{})
	}

	actions := make([]gomcts.Action, len(candidates))
	for i := 0; i < len(candidates); i++ {
		actions[i] = &RedTenGameAction{hand: candidates[i], nextToMove: t.nextToMove}
	}
	return actions
}

func (t *RedTenGameState) IsGameEnded() bool {
	_, ended := t.EvaluateGame()
	return ended
}

func (t *RedTenGameState) NextToMove() int8 {
	return t.nextToMove
}

// CreateRedTenInitialGameState - initializes tic tac toe game state
func CreateRedTenInitialGameState(a, b, preHand []int) RedTenGameState {
	state := RedTenGameState{a: a, b: b, preHand: preHand, nextToMove: 1}
	return state
}

// ///////////////////////////////////////////////////////////////////////

const (
	WIN_REDTEN     = 1
	WIN_REDTEN_NON = 2
)

var (
	playerNum = 5
)

//5人红十
// RedTen5GameAction - action on a tic tac toe board game
type RedTen5GameAction struct {
	hands [][]int //自己或者多个下家出牌
}

func (t *RedTen5GameAction) Hash() uint64 {
	return murmur3.Sum64(structhash.Dump(t, 1))
}

func (t *RedTen5GameAction) ApplyTo(s gomcts.GameState) gomcts.GameState {
	var turn int
	state := s.(*RedTen5GameState)

	newState := &RedTen5GameState{
		myRole:  state.myRole,
		inHands: make([][]int, playerNum),
		//preHand:    t.hands[len(t.hands)-1],
		nextToMove: -1 * state.nextToMove,
	}

	if state.nextToMove == 1 {
		//轮到我
		turn = state.myRole
		if len(t.hands[0]) == 0 {
			//我跳过
			newState.preHand = state.preHand
			newState.preRole = state.preRole
			newState.inHands[turn] = state.inHands[turn]
		} else {
			//我出
			newState.preHand = t.hands[len(t.hands)-1]
			newState.preRole = state.myRole
			newState.inHands[turn] = card.RemoveCards2(state.inHands[turn], t.hands[0])
		}

		for i := 0; i < playerNum; i++ {
			turn = (turn + 1) % playerNum
			newState.inHands[turn] = state.inHands[turn]
		}
	} else {
		//我的牌不变
		turn = state.myRole
		newState.inHands[turn] = state.inHands[turn]
		for i := 0; i < len(t.hands); i++ {
			turn = (state.myRole + 1) % playerNum
			if len(t.hands[i]) > 0 {
				newState.preHand = t.hands[i]
				newState.preRole = turn
				newState.inHands[turn] = card.RemoveCards2(state.inHands[turn], t.hands[i])
			} else {
				newState.inHands[turn] = state.inHands[turn]
			}
		}

		if newState.preHand == nil {
			//都跳过了
			newState.preHand = state.preHand
			newState.preRole = state.preRole
		}
	}

	//fmt.Printf("\nApplyTo,state:%v,newState:%v,action:%v", state, newState, t)
	return newState
}

type RedTen5GameState struct {
	myRole     int
	inHands    [][]int
	preHand    []int
	preRole    int
	nextToMove int8
}

func (t *RedTen5GameState) Hash() uint64 {
	return murmur3.Sum64(structhash.Dump(t, 1))
}

func (t *RedTen5GameState) EvaluateGame() (gomcts.GameResult, bool) {
	return 0, true
}

func (t *RedTen5GameState) generateActions(turn int, n int) []gomcts.Action {
	var candidates [][]int
	var actions [][][]int
	for i := 0; i < n; i++ {
		if t.preHand == nil || len(t.preHand) == 0 || t.preRole == turn {
			//新出或者上轮此人出没人要
			candidates = aviableCandidates(t.inHands[turn])
		} else {
			candidates = aviableBiggerCandidates(t.inHands[turn], t.preHand)
			candidates = append(candidates, []int{})
		}

		tmpActions := actions
		if tmpActions == nil {
			actions = make([][][]int, 0, len(candidates))
			for j := 0; j < len(candidates); j++ {
				action := make([][]int, n)
				action[i] = candidates[j]
				actions = append(actions, action)
			}
		} else {
			actions = make([][][]int, 0, len(candidates)*len(tmpActions))
			for j := 0; j < len(candidates); j++ {
				for k := 0; k < len(tmpActions); k++ {
					action := make([][]int, n)
					copy(action, tmpActions[k])
					action[i] = candidates[j]
					actions = append(actions, action)
				}
			}
		}

		turn = (t.myRole + 1) % playerNum
	}

	iActions := make([]gomcts.Action, len(actions))
	for i := 0; i < len(actions); i++ {
		iActions[i] = &RedTen5GameAction{hands: actions[i]}
	}

	return iActions
}

func (t *RedTen5GameState) GetLegalActions() []gomcts.Action {
	var actions []gomcts.Action

	if t.nextToMove == 1 {
		//轮到我了
		actions = t.generateActions(t.myRole, 1)
	} else {
		actions = t.generateActions((t.myRole+1)%playerNum, playerNum-1)
	}

	return actions
}

func (t *RedTen5GameState) IsGameEnded() bool {
	_, ended := t.EvaluateGame()
	return ended
}

func (t *RedTen5GameState) NextToMove() int8 {
	return t.nextToMove
}

// CreateRedTen5GameState - initializes game state
func CreateRedTen5GameState(myRole int, inHands [][]int, preHand []int, preRole int) RedTen5GameState {
	state := RedTen5GameState{myRole: myRole, inHands: inHands, preHand: preHand, preRole: preRole, nextToMove: 1}
	return state
}
