package engine

import "testing"

func TestPlayCard(t *testing.T) {
	game := NewGame("../data/combat_card.json")
	game.PlayerA.Draw()
	game.PlayCard(0, true)
}
