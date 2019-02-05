package sm

// #cgo CFLAGS: -DWHATEVER_YOU_WANT_TO_INDICATE_CGO=1
// #cgo CXXFLAGS: -DWHATEVER_YOU_WANT_TO_INDICATE_CGO=1
// #include "sm.h"
import "C"

import (
	"errors"
	"unsafe"
)

func init() {
	C.kusokurae_global_init()
}

// Enum: kusokurae_game_status_t
const (
	StatusNull int32 = iota
	StatusInit
	StatusPlay
	StatusFinish

	StatusMax
)

// Enum: kusokurae_card_suit_t
const (
	SuitXiang   int32 = -1
	SuitYoutiao       = 0
	SuitBaozi         = 1
	SuitOther         = 2
)

// Enum: kusokurae_round_status_t
const (
	RoundWaiting int32 = iota
	RoundActive
	RoundDone
)

// Errors from underlying library.
// Don't forget also to change here after adding new error codes in C interface.
var (
	ErrNullPtr       = errors.New("KUSOKURAE_ERROR_NULLPTR")
	ErrBadNPlayers   = errors.New("KUSOKURAE_ERROR_BAD_NUMBER_OF_PLAYERS")
	ErrUninitialized = errors.New("KUSOKURAE_ERROR_UNINITIALIZED")

	ErrUnknown = errors.New("Unknown")
)

var errMap = map[C.kusokurae_error_t]error{
	C.KUSOKURAE_SUCCESS:                     nil,
	C.KUSOKURAE_ERROR_NULLPTR:               ErrNullPtr,
	C.KUSOKURAE_ERROR_BAD_NUMBER_OF_PLAYERS: ErrBadNPlayers,
	C.KUSOKURAE_ERROR_UNINITIALIZED:         ErrUninitialized,
}

// GameConfig has the same memory layout with C.kusokurae_game_config_t.
type GameConfig struct {
	NumPlayers int32
}

// Card has the same memory layout with C.kusokurae_card_t.
type Card struct {
	displayOrder uint32
	Suit         int32 // C.kusokurae_card_suit_t
	Rank         int32
}

// Player has the same memory layout with C.kusokurae_player_t.
type Player struct {
	index      int32
	active     int32 // C.kusokurae_round_status_t
	hand       [C.KUSOKURAE_MAX_HAND_CARDS]Card
	numCards   int32
	cardsTaken int32
	score      int32
}

// GameState has the same memory layout with C.kusokurae_game_state_t.
type GameState struct {
	cfg         GameConfig
	status      int32
	players     [C.KUSOKURAE_MAX_PLAYERS]Player
	numRound    int32
	ghostHolder int32
	curRound    [C.KUSOKURAE_MAX_PLAYERS]Card
}

// RoundState corresponds to C.kusokurae_round_state_t, but does not preserve
// its memory layout - a bit of conversion needs to be done when providing this
// to user Go code.
type RoundState struct {
	Seq          int
	IsDoubled    bool
	ScoreOnBoard int
	Leader       *Player
}

func errcode2Go(code C.kusokurae_error_t) (err error) {
	err, ok := errMap[code]
	if !ok {
		err = ErrUnknown
	}
	return
}

// GetHandCards returns a slice holding the player's cards. It operates in
// constant time.
func (p *Player) GetHandCards() []Card {
	return p.hand[0:p.numCards]
}

// NewGame creates a new game state with specified number of players.
func NewGame(cfg GameConfig) (ret *GameState, err error) {
	ret = &GameState{}
	pret := unsafe.Pointer(ret)
	pcfg := unsafe.Pointer(&cfg)
	err = errcode2Go(C.kusokurae_game_init(
		(*C.kusokurae_game_state_t)(pret),
		(*C.kusokurae_game_config_t)(pcfg),
	))
	return
}

// Start deals cards to each player and begins waiting for play from the first
// player.
func (g *GameState) Start() (err error) {
	pg := unsafe.Pointer(g)
	err = errcode2Go(C.kusokurae_game_start((*C.kusokurae_game_state_t)(pg)))
	return
}

// IsFinalRound checks if the game is in (or after) its last round.
func (g *GameState) IsFinalRound() bool {
	for i := range g.players {
		if g.players[i].numCards > 1 {
			return false
		}
	}
	return true
}
