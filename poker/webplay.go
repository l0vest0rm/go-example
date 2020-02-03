package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

const (
	REQ_CHEAT = 1
	RSP_CHEAT = 2

	REQ_CREATE_AI_TABLE = 3 //创建一个机器人局
	RSP_CREATE_AI_TABLE = 4

	REQ_LOGIN = 11
	RSP_LOGIN = 12

	REQ_ROOM_LIST = 13
	RSP_ROOM_LIST = 14

	REQ_TABLE_LIST = 15
	RSP_TABLE_LIST = 16

	REQ_JOIN_ROOM = 17
	RSP_JOIN_ROOM = 18

	REQ_JOIN_TABLE = 19
	RSP_JOIN_TABLE = 20

	REQ_NEW_TABLE = 21
	RSP_NEW_TABLE = 22

	YOUR_TURN      = 23 //轮到你出牌
	YOU_SHOT       = 24 //你可以压
	INVALID_POCKER = 25 //牌型不合法

	REQ_DEAL_POKER = 31
	RSP_DEAL_POKER = 32

	REQ_CALL_SCORE = 33
	RSP_CALL_SCORE = 34

	REQ_SHOW_POKER = 35
	RSP_SHOW_POKER = 36

	REQ_SHOT_POKER = 37
	RSP_SHOT_POKER = 38

	REQ_GAME_OVER = 41
	RSP_GAME_OVER = 42

	REQ_CHAT = 43
	RSP_CHAT = 44

	REQ_RESTART = 45
	RSP_RESTART = 46
)

type Message struct {
	Code    int         `json:"code"`
	Uid     int         `json:"uid"`
	TableId int         `json:"tableId,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PlayerInfo struct {
	Uid  int    `json:"uid"`
	Name string `json:"name"`
}

type TableInfo struct {
	TableId int `json:"tableId"`
	Players int `json:"players"`
}

type UserInfo struct {
	uid      int
	tableId  int
	playerId int
}

type Games struct {
	mu          sync.RWMutex
	nextTableId int
	gamesMap    map[int]IGame
}

var (
	games Games
)

func webPlay() {
	isDebug = true
	games = Games{
		nextTableId: 101,
		gamesMap:    make(map[int]IGame),
	}

	http.Handle("/ws", websocket.Handler(webSocketRcv))
	http.HandleFunc("/", servStatic)
	http.ListenAndServe(":8080", nil)
}

func webSocketRcv(ws *websocket.Conn) {
	var err error
	user := UserInfo{}
	user.uid = 1710

	for {
		var req []byte
		var msg Message

		if err = websocket.Message.Receive(ws, &req); err != nil {
			fmt.Println(err)
			ws.Close()
			return
		}

		err := json.Unmarshal(req, &msg)
		if err != nil {
			fmt.Println("json.Unmarshal err:", err)
		}

		switch msg.Code {
		case REQ_CREATE_AI_TABLE:
			user.uid = msg.Uid
			onReqCreateAiTable(ws, &user)
		case REQ_TABLE_LIST:
			user.uid = msg.Uid
			onReqTableList(ws, &user)
		case REQ_DEAL_POKER:
			onReqDealPoker(ws, &user)
		case REQ_SHOT_POKER:
			onReqShotPoker(ws, &user, &msg)
		default:
			fmt.Printf("unknown req:%d\n", msg.Code)
			continue
		}
	}
}

func gameOver(ws *websocket.Conn, game IGame, uid int) {
	scores := game.CalcScores()
	rsp := Message{
		Code: RSP_GAME_OVER,
		Uid:  uid,
		Data: scores,
	}
	webSocketSendMsg(ws, rsp)
}

func onReqCreateAiTable(ws *websocket.Conn, user *UserInfo) {
	tableId := prepreNewGame()
	user.tableId = tableId

	data := []PlayerInfo{
		{Uid: user.uid, Name: "terry"},
		{Uid: 1, Name: "player1"},
		{Uid: 2, Name: "player2"},
		{Uid: 3, Name: "player3"},
		{Uid: 4, Name: "player4"},
	}

	rsp := Message{
		Code:    RSP_CREATE_AI_TABLE,
		Uid:     user.uid,
		TableId: user.tableId,
		Data:    data,
	}

	webSocketSendMsg(ws, rsp)
}

func onReqTableList(ws *websocket.Conn, user *UserInfo) {

	data := []TableInfo{
		{TableId: 1, Players: 0},
		{TableId: 2, Players: 0},
	}

	rsp := Message{
		Code: RSP_TABLE_LIST,
		Uid:  user.uid,
		Data: data,
	}

	webSocketSendMsg(ws, rsp)
}

func onReqDealPoker(ws *websocket.Conn, user *UserInfo) {
	games.mu.RLock()
	game, ok := games.gamesMap[user.tableId]
	games.mu.RUnlock()
	if !ok {
		fmt.Printf("table not ready:%d,uid:%d", user.tableId, user.uid)
		return
	}

	data := []PlayerInfo{
		{Uid: user.uid, Name: "me"},
		{Uid: 1, Name: "player1"},
		{Uid: 2, Name: "player2"},
		{Uid: 3, Name: "player3"},
		{Uid: 4, Name: "player4"},
	}

	rsp := Message{
		Code: RSP_JOIN_TABLE,
		Uid:  user.uid,
		Data: data,
	}

	webSocketSendMsg(ws, rsp)

	remainCards := game.RemainCards(0)
	rsp = Message{
		Code: RSP_DEAL_POKER,
		Uid:  user.uid,
		Data: remainCards,
	}

	//发牌给玩家
	webSocketSendMsg(ws, rsp)

	//开始出牌
	prePlayerId := -1
	nextRound(ws, game, user.uid, prePlayerId)

}

func onReqShotPoker(ws *websocket.Conn, user *UserInfo, msg *Message) {
	game := getGame(user.tableId)
	if game == nil {
		fmt.Printf("table not found:%d,uid:%d", user.tableId, user.uid)
		return
	}

	cardsFls := msg.Data.([]interface{})
	cards := make([]int, 0)
	valid := true
	for _, card := range cardsFls {
		cards = append(cards, int(card.(float64)))
	}

	cards, valid = game.PlayerHand(0, cards)
	if cards == nil {
		cards = []int{}
	}

	var rsp Message
	if valid {
		rsp = Message{
			Code: RSP_SHOT_POKER,
			Uid:  user.uid,
			Data: cards,
		}
	} else {
		rsp = Message{
			Code: INVALID_POCKER,
			Uid:  user.uid,
		}
	}

	webSocketSendMsg(ws, rsp)
	if !valid {
		return
	}

	nextRound(ws, game, user.uid, 0)
}

func nextRound(ws *websocket.Conn, game IGame, uid int, prePlayerId int) {
	var rsp Message
	var code int
	for {
		playerId := game.NextPlayer(prePlayerId)
		if playerId == -1 {
			Println("\ngame over")
			gameOver(ws, game, uid)
			return
		}
		prePlayerId = playerId
		if playerId == 0 {
			if game.isYourTurn(playerId) {
				code = YOUR_TURN
			} else {
				code = YOU_SHOT
			}
			rsp = Message{
				Code: code,
				Uid:  uid,
			}

			webSocketSendMsg(ws, rsp)
			break
		}

		time.Sleep(time.Second * 2)

		cards, _ := game.PlayerHand(playerId, nil)
		if cards == nil {
			cards = []int{}
		}
		rsp = Message{
			Code: RSP_SHOT_POKER,
			Uid:  playerId,
			Data: cards,
		}

		webSocketSendMsg(ws, rsp)
	}
}

func prepreNewGame() int {
	players := []int{0, 2, 2, 2, 2}
	game := NewRedTen(players)

	games.mu.Lock()
	tableId := games.nextTableId
	games.nextTableId += 1
	games.gamesMap[tableId] = game
	games.mu.Unlock()

	return tableId
}

func webSocketSendMsg(ws *websocket.Conn, msg interface{}) {
	bs, err := json.Marshal(msg)
	if err != nil {
		panic("webSocketSendMsg json.Marshal fail")
	}

	//发送给客户端
	if err := websocket.Message.Send(ws, string(bs)); err != nil {
		fmt.Println(err)
	}
}

func getGame(tableId int) IGame {
	games.mu.RLock()
	game, ok := games.gamesMap[tableId]
	games.mu.RUnlock()
	if !ok {
		fmt.Printf("table not found:%d", tableId)
		return nil
	}

	return game
}
