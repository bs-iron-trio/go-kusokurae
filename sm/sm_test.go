package sm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	_, err := NewGame(GameConfig{
		NumPlayers: 2,
	}, nil)
	assert.Equal(t, ErrBadNPlayers, err)

	correctCfg := GameConfig{
		NumPlayers: 3,
	}
	state, err := NewGame(correctCfg, nil)
	assert.NoError(t, err)
	assert.Equal(t, correctCfg, state.cfg)
	assert.Equal(t, 1, int(state.players[0].index))
	assert.Equal(t, 2, int(state.players[1].index))
	assert.Equal(t, 3, int(state.players[2].index))
	assert.Equal(t, 0, int(state.players[3].index))
}

func TestActivePlayerNil(t *testing.T) {
	state, err := NewGame(GameConfig{
		NumPlayers: 3,
	}, nil)
	assert.NoError(t, err)

	// bs@bs-newnb-w10:~/go/src/github.com/bs-iron-trio/go-kusokurae/sm$ go test
	// --- FAIL: TestActivePlayerNil (0.00s)
	// 		sm_test.go:34:
	// 						Error Trace:    sm_test.go:34
	// 						Error:          Not equal:
	// 										expected: <nil>(<nil>)
	// 										actual  : *sm.Player((*sm.Player)(nil))
	// 						Test:           TestActivePlayerNil
	// FAIL
	// exit status 1
	// FAIL    github.com/bs-iron-trio/go-kusokurae/sm 0.034s
	//assert.Equal(t, nil, state.GetActivePlayer())
	//assert.EqualValues(t, nil, state.GetActivePlayer())
	assert.Nil(t, state.GetActivePlayer())
}

func TestStateCB(t *testing.T) {
	var calls int
	var recordedNewState GameStatus
	state, err := NewGame(GameConfig{
		NumPlayers: 4,
	}, func(newState GameStatus) {
		calls++
		recordedNewState = newState
	})
	assert.NoError(t, err)

	err = state.Start()
	assert.NoError(t, err)
	assert.Equal(t, 1, calls)
	assert.Equal(t, StatusPlay, recordedNewState)
}

func TestGameStart(t *testing.T) {
	state, err := NewGame(GameConfig{
		NumPlayers: 3,
	}, nil)
	assert.NoError(t, err)

	err = state.Start()
	assert.NoError(t, err)
	assert.Equal(t, RoundActive, state.players[0].active)
	assert.Equal(t, RoundWaiting, state.players[1].active)
	assert.Equal(t, RoundWaiting, state.players[2].active)
	assert.Equal(t, StatusPlay, state.status)
	assert.Equal(t, &state.players[0], state.GetActivePlayer())

	t.Log(state.players[0].allCards)
	t.Log(state.players[1].allCards)
	t.Log(state.players[2].allCards)
	// Verify dealing: ensure no duplicate card
	cards := make(map[uint32]bool)
	var i, j int
	var order uint32
	for i = 0; i < 3; i++ {
		for j = 0; j < 11; j++ {
			order = state.players[i].allCards[j].displayOrder
			if cards[order] {
				t.Errorf("Duplicate card %+v", state.players[i].allCards[j])
			}
			cards[order] = true
		}
	}
}

func TestCardString(t *testing.T) {
	assert.Equal(t, "8(-1)", fmt.Sprintf("%v", Card{
		suit: SuitXiang,
		rank: 8,
	}))
	assert.Equal(t, "10(x2)", fmt.Sprintf("%v", Card{
		displayOrder: 3,
		suit:         SuitOther,
		rank:         10,
	}))
	assert.Equal(t, "1(1),played=1", Card{
		suit:  SuitBaozi,
		rank:  1,
		flags: 1,
	}.String())
	assert.Equal(t, "[2(0) 3(0)]", fmt.Sprint([]Card{
		{0, SuitYoutiao, 2, 128},
		{0, SuitYoutiao, 3, 128},
	}))
}
