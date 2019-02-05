package sm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	_, err := NewGame(GameConfig{
		NumPlayers: 2,
	})
	assert.Equal(t, ErrBadNPlayers, err)

	correctCfg := GameConfig{
		NumPlayers: 3,
	}
	state, err := NewGame(correctCfg)
	assert.NoError(t, err)
	assert.Equal(t, correctCfg, state.cfg)
	assert.Equal(t, StatusInit, state.status)
	assert.Equal(t, 1, int(state.players[0].index))
	assert.Equal(t, 2, int(state.players[1].index))
	assert.Equal(t, 3, int(state.players[2].index))
	assert.Equal(t, 0, int(state.players[3].index))
}

func TestGameStart(t *testing.T) {
	state, err := NewGame(GameConfig{
		NumPlayers: 3,
	})
	assert.NoError(t, err)

	err = state.Start()
	assert.NoError(t, err)
	assert.Equal(t, RoundActive, state.players[0].active)
	assert.Equal(t, RoundWaiting, state.players[1].active)
	assert.Equal(t, RoundWaiting, state.players[2].active)
	assert.Equal(t, StatusPlay, state.status)

	t.Log(state.players[0].hand)
	t.Log(state.players[1].hand)
	t.Log(state.players[2].hand)
	cards := make(map[uint32]bool)
	var i, j int
	var order uint32
	for i = 0; i < 3; i++ {
		for j = 0; j < 11; j++ {
			order = state.players[i].hand[j].displayOrder
			if cards[order] {
				t.Errorf("Duplicate card %+v", state.players[i].hand[j])
			}
			cards[order] = true
		}
	}
}
