package doudizhu

import (
	"../card"
	"../gomcts"
)

//角色
const (
	ROLE_DIZHU = 0 //地主
	ROLE_NON1  = 1 //农名1
	ROLE_NONG2 = 2 //农名2
)

//谁赢了
const (
	WIN_DIZHU = 1
	WIN_NONG  = 2
)

// RedTenGameAction - action on a tic tac toe board game
type DoudizhuGameAction struct {
	hands [][]int //自己或者多个下家出牌
}

func (t *DoudizhuGameAction) ApplyTo(s gomcts.GameState) gomcts.GameState {
	var err error
	var turn int
	state := s.(*DoudizhuGameState)

	newState := &DoudizhuGameState{
		myRole:  state.myRole,
		inHands: make([][]int, 3),
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
			newState.inHands[turn], err = card.RemoveCards(state.inHands[turn], t.hands[0])
			if err != nil {
				panic(err)
			}
		}

		turn = (turn + 1) % 3
		newState.inHands[turn] = state.inHands[turn]
		turn = (turn + 1) % 3
		newState.inHands[turn] = state.inHands[turn]
	} else {
		turn = (state.myRole + 1) % 3
		if len(t.hands[0]) > 0 {
			newState.preHand = t.hands[0]
			newState.preRole = turn
			newState.inHands[turn], err = card.RemoveCards(state.inHands[turn], t.hands[0])
			if err != nil {
				panic(err)
			}
		} else {
			newState.inHands[turn] = state.inHands[turn]
		}

		turn = (turn + 1) % 3
		if len(t.hands[1]) > 0 {
			newState.preHand = t.hands[1]
			newState.preRole = turn
			newState.inHands[turn], err = card.RemoveCards(state.inHands[turn], t.hands[1])
			if err != nil {
				panic(err)
			}
		} else {
			newState.inHands[turn] = state.inHands[turn]
		}

		turn = (turn + 1) % 3
		newState.inHands[turn] = state.inHands[turn]

		if newState.preHand == nil {
			//都跳过了
			newState.preHand = state.preHand
			newState.preRole = state.preRole
		}
	}

	return newState
}

type DoudizhuGameState struct {
	myRole     int
	inHands    [][]int
	preHand    []int
	preRole    int
	nextToMove int8
}

func (t *DoudizhuGameState) EvaluateGame() (gomcts.GameResult, bool) {
	var whoWin int
	if len(t.inHands[ROLE_DIZHU]) == 0 {
		whoWin = WIN_DIZHU
	} else if len(t.inHands[ROLE_NON1]) == 0 || len(t.inHands[ROLE_NON2]) == 0 {
		whoWin = WIN_NONG
	}
	if t.myRole == ROLE_DIZHU {
		if whoWin == WIN_DIZHU {
			return 1, true
		} else if whoWin == WIN_NONG {
			return -1, true
		} else {
			return 0, false
		}
	} else {
		if whoWin == WIN_DIZHU {
			return -1, true
		} else if whoWin == WIN_NONG {
			return 1, true
		} else {
			return 0, false
		}
	}
}

func (t *DoudizhuGameState) GetLegalActions() []gomcts.Action {
	var turn int
	var candidates [][]int
	var candidates2 [][]int
	var actions []gomcts.Action

	if t.nextToMove == 1 {
		//轮到我了
		actions = make([]gomcts.Action, len(candidates))
		if t.preHand == nil || len(t.preHand) == 0 || t.preRole == t.myRole {
			//新出或者上轮我出的没人要
			candidates = aviableCandidates(t.inHands[t.myRole])
		} else {
			candidates = aviableBiggerCandidates(t.inHands[t.myRole], t.preHand)
			candidates = append(candidates, []int{})
		}

		for i := 0; i < len(candidates); i++ {
			actions[i] = &DoudizhuGameAction{hands: make([][]int, 1)}
			actions[i].(*DoudizhuGameAction).hands[0] = candidates[i]
		}

	} else {
		actions = make([]gomcts.Action, 0)
		turn = (t.myRole + 1) % 3
		if t.preHand == nil || len(t.preHand) == 0 || t.preRole == turn {
			//新出或者上轮此人出没人要
			candidates = aviableCandidates(t.inHands[turn])
		} else {
			candidates = aviableBiggerCandidates(t.inHands[turn], t.preHand)
		}

		turn = (turn + 1) % 3
		for i := 0; i < len(candidates); i++ {
			candidates2 = aviableBiggerCandidates(t.inHands[turn], candidates[i])
			for j := 0; j < len(candidates2); j++ {
				action := &DoudizhuGameAction{hands: make([][]int, 2)}
				action.hands[0] = candidates[i]
				action.hands[1] = candidates2[j]
				actions = append(actions, action)
			}
		}

		//考虑前一个人不出的情况
		if t.preRole == turn {
			//上轮此人出没人要
			candidates2 = aviableCandidates(t.inHands[turn])
		} else {
			candidates2 = aviableBiggerCandidates(t.inHands[turn], t.preHand)
		}
		for j := 0; j < len(candidates2); j++ {
			action := &DoudizhuGameAction{hands: make([][]int, 2)}
			action.hands[0] = []int{}
			action.hands[1] = candidates2[j]
			actions = append(actions, action)
		}
	}

	return actions
}

func (t *DoudizhuGameState) IsGameEnded() bool {
	_, ended := t.EvaluateGame()
	return ended
}

func (t *DoudizhuGameState) NextToMove() int8 {
	return t.nextToMove
}

// CreateDoudizhuGameState - initializes game state
func CreateDoudizhuGameState(myRole int, inHands [][]int, preHand []int, preRole int) DoudizhuGameState {
	state := DoudizhuGameState{myRole: myRole, inHands: inHands, preHand: preHand, preRole: preRole, nextToMove: 1}
	return state
}
