package gomcts

// GameResult - number representing a game result
type GameResult float64

// Action - interface representing entity that can be applied to a game state (generating the next game state)
type Action interface {
	Hash() uint64
	ApplyTo(GameState) GameState
}

// GameState - state of the game interface
type GameState interface {
	Hash() uint64
	EvaluateGame() (GameResult, bool)
	GetLegalActions() []Action
	IsGameEnded() bool
	NextToMove() int8
}
