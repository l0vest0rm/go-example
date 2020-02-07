package doudizhu

import (
	"../gomcts"
)

type Role int

//角色
const (
	ROLE_DIZHU = 0 //地主
	ROLE_NON1  = 1 //农名1
	ROLE_NONG2 = 2 //农名2
)

// RedTenGameAction - action on a tic tac toe board game
type DoudizhuGameAction struct {
	hand1      []int //自己出时只有hand1
	hand2      []int //如果是对方出时下下一个人是hand2
	nextToMove int8
}

func (t *DoudizhuGameAction) ApplyTo(s gomcts.GameState) gomcts.GameState {
	var err error
	state := s.(*DoudizhuGameAction)
	if state.nextToMove != t.nextToMove {
		panic("wrong turn")
	}

	newState := &DoudizhuGameAction{
		a:          state.a,
		b:          state.b,
		preHand:    t.hand,
		nextToMove: -1 * state.nextToMove,
	}

	if t.nextToMove == 1 {
		newState.a, err = removeCards(newState.a, t.hand)
	} else {
		newState.b, err = removeCards(newState.b, t.hand)
	}

	if err != nil {
		panic(err)
	}

	return newState
}

type DoudizhuGameState struct {
	a          []int
	b          []int
	c          []int
	preHand    []int
	nextToMove int8
}

func (t *DoudizhuGameState) EvaluateGame() (gomcts.GameResult, bool) {
	if len(t.a) == 0 {
		return 1, true
	} else if len(t.b) == 0 {
		return -1, true
	} else {
		return 0, false
	}
}

func (t *DoudizhuGameState) GetLegalActions() []gomcts.Action {
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
		actions[i] = &DoudizhuGameState{hand: candidates[i], nextToMove: t.nextToMove}
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
func CreateDoudizhuGameState(role Role, a, b, c, preHand []int) DoudizhuGameState {
	state := DoudizhuGameState{a: a, b: b, preHand: preHand, nextToMove: 1}
	return state
}
