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

func TestPickCombo(t *testing.T) {
	game := NewGame("../data/combat_card.json")
	if err := game.PlayCard(0, true); err != nil { // Player A plays a card
		t.Error(err)
	}
	if err := game.PlayCard(0, true); err != nil { // Player B plays a card
		t.Error(err)
	}
	if err := game.Abandon(); err != nil { // Player A abandons
		t.Error(err)
	}
	if err := game.PickCombo(0); err != nil { // Player B picks combos
		t.Error(err)
	}
	if err := game.SetCurrentPlayer(true); err != nil {
		t.Error(err)
	}
	if game.State != Combat {
		t.Errorf("New game round has not began, game state is %s", game.State)
	}

	numberOfCardsPlayerA := len(game.PlayerA.Deck) + len(game.PlayerA.Hand)
	if numberOfCardsPlayerA != 8 {
		t.Errorf("Total number of cards in players deck and hand is %d, expected 8", numberOfCardsPlayerA)
	}
	numberOfCardsPlayerB := len(game.PlayerB.Deck) + len(game.PlayerB.Hand)
	if numberOfCardsPlayerB != 10 {
		t.Errorf("Total number of cards in players deck and hand is %d, expected 10", numberOfCardsPlayerB)
	}
}
