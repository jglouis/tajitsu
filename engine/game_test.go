package engine

import "testing"

func TestPlayCard(t *testing.T) {
	game := NewGame("../data/combat_card.json")
	game.PlayerA.Draw()
	game.PlayCard(0, true)
}

func TestAbandonToVictory(t *testing.T) {
	game := NewGame("../data/combat_card.json")
	if err := game.Abandon(); err != nil {
		t.Error(err)
	}
	if err := game.SetCurrentPlayer(true); err != nil {
		t.Error(err)
	}
	if err := game.Abandon(); err != nil {
		t.Error(err)
	}

	if game.ScoreB != 2 {
		t.Errorf("Player B score is %d, expected 2", game.ScoreB)
	}
	if game.ScoreA != 0 {
		t.Errorf("Player A score is %d, expected 0", game.ScoreB)
	}
	if game.State != End {
		t.Errorf("Game has not ended, game state is %s", game.State)
	}
}
