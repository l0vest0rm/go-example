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

	START_PLAY = 23

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
	TableId int         `json:"tableId"`
	Uid     int         `json:"uid"`
	Data    interface{} `json:"data"`
}

type PlayerInfo struct {
	Uid  int    `json:"uid"`
	Name string `json:"name"`
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
	games = Games{
		nextTableId: 1,
		gamesMap:    make(map[int]IGame),
	}

	http.Handle("/ws", websocket.Handler(webSocketRcv))
	http.HandleFunc("/", servStatic)
	http.ListenAndServe(":8080", nil)
}

func webSocketRcv(ws *websocket.Conn) {
	var err error
	for {
		var req []byte
		var msg Message

		if err = websocket.Message.Receive(ws, &req); err != nil {
			panic(err)
			continue
		}

		err := json.Unmarshal(req, &msg)
		if err != nil {
			fmt.Println("json.Unmarshal err:", err)
		}

		switch msg.Code {
		case REQ_JOIN_ROOM:
			onReqJoinRoom(ws, &msg)
		case REQ_JOIN_TABLE:
			onReqJoinTable(ws, &msg)
		case REQ_SHOT_POKER:
			onReqShotPoker(ws, &msg)
		default:
			fmt.Printf("unknown req:%d\n", msg.Code)
			continue
		}
	}
}

func onReqJoinRoom(ws *websocket.Conn, msg *Message) {
	msg.Code = RSP_JOIN_ROOM
	webSocketSendMsg(ws, msg)
}

func onReqJoinTable(ws *websocket.Conn, msg *Message) {
	tableId := msg.TableId
	if tableId < 1 {
		tableId = prepreGame()
	}

	games.mu.RLock()
	game, ok := games.gamesMap[tableId]
	games.mu.RUnlock()
	if !ok {
		fmt.Printf("table not ready:%d,uid:%d", tableId, msg.Uid)
		return
	}

	data := []PlayerInfo{
		{Uid: msg.Uid, Name: "terry"},
		{Uid: 1, Name: "player1"},
		{Uid: 2, Name: "player2"},
		{Uid: 3, Name: "player3"},
		{Uid: 4, Name: "player4"},
	}

	rsp := Message{
		Code:    RSP_JOIN_TABLE,
		TableId: tableId,
		Uid:     msg.Uid,
		Data:    data,
	}

	webSocketSendMsg(ws, rsp)

	remainCards := game.RemainCards(0)
	rsp = Message{
		Code:    RSP_DEAL_POKER,
		TableId: tableId,
		Uid:     msg.Uid,
		Data:    remainCards,
	}

	//发牌给玩家
	webSocketSendMsg(ws, rsp)

	//开始出牌
	prePlayerId := -1
	nextRound(ws, tableId, msg.Uid, prePlayerId)

}

func onReqShotPoker(ws *websocket.Conn, msg *Message) {
	game := getGame(msg.TableId)
	if game == nil {
		fmt.Printf("table not found:%d,uid:%d", msg.TableId, msg.Uid)
		return
	}

	cardsFls := msg.Data.([]interface{})
	cards := make([]int, 0)
	for _, card := range cardsFls {
		cards = append(cards, int(card.(float64)))
	}

	if len(cards) > 0 {
		game.RecordHand(0, cards)
	}

	rsp := Message{
		Code:    RSP_SHOT_POKER,
		TableId: msg.TableId,
		Uid:     msg.Uid,
		Data:    cards,
	}

	webSocketSendMsg(ws, rsp)

	nextRound(ws, msg.TableId, msg.Uid, 0)

}

func nextRound(ws *websocket.Conn, tableId int, uid int, prePlayerId int) {
	game := getGame(tableId)
	if game == nil {
		fmt.Printf("table not found:%d,uid:%d", tableId, uid)
		return
	}

	var rsp Message
	for {
		time.Sleep(time.Second * 2)
		playerId := game.NextPlayer(prePlayerId)
		prePlayerId = playerId
		if playerId == 0 {
			rsp = Message{
				Code:    START_PLAY,
				TableId: tableId,
				Uid:     uid,
			}

			webSocketSendMsg(ws, rsp)
			break
		}

		cards := game.PlayerHand(playerId)
		if cards == nil {
			cards = []int{}
		}
		rsp = Message{
			Code:    RSP_SHOT_POKER,
			TableId: tableId,
			Uid:     playerId,
			Data:    cards,
		}

		webSocketSendMsg(ws, rsp)
	}
}

func prepreGame() int {
	players := []int{0, 1, 1, 1, 1}
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
