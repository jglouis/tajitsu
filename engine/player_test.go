package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var combatCards []*CombatCard

func init() {
	f, e := ioutil.ReadFile("../data/combat_card.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	json.Unmarshal(f, &combatCards)
}

func setUpPlayer() Player {
	var player Player
	for _, combatCard := range combatCards {
		player.Deck = append(player.Deck, combatCard)
	}
	return player
}

func TestDraw(t *testing.T) {
	player := setUpPlayer()

	for i := 0; i < 10; i++ {
		card := player.Draw()
		fmt.Println(card)
	}
	if len(player.Deck) != 0 {
		t.Errorf("Expected player's deck to be empty, but there are %d cards left", len(player.Deck))
	}
	if len(player.Hand) != 9 {
		t.Errorf("Expected player's hand to have nine cards, but has %d", len(player.Hand))
	}
}

func TestDiscard(t *testing.T) {
	player := setUpPlayer()
	player.Discard(0) // Test discarding a empty hand
	for i := 0; i < 9; i++ {
		player.Draw()
		player.Discard(0)
	}
	if len(player.Deck) != 0 {
		t.Errorf("Expected player's deck to be empty, but there are %d cards left", len(player.Deck))
	}
	if len(player.Hand) != 0 {
		t.Errorf("Expected player's hand to be empty, but has %d", len(player.Hand))
	}
	if len(player.DiscardPile) != 9 {
		t.Errorf("Expected player's discard pile to have nine cards, but has %d", len(player.DiscardPile))
	}
}

func TestDeckShuffle(t *testing.T) {
	player := setUpPlayer()
	player.DeckShuffle()
	if len(player.Deck) != 9 {
		t.Errorf("Expected player's deck tohave nine cards, but there are %d cards left", len(player.Deck))
	}
	if len(player.Hand) != 0 {
		t.Errorf("Expected player's hand to be empty, but has %d", len(player.Hand))
	}
	if len(player.DiscardPile) != 0 {
		t.Errorf("Expected player's discard pile to be empty, but has %d", len(player.DiscardPile))
	}
}

func TestPlayCard(t *testing.T) {
	game := NewGame("../data/combat_card.json")
	game.PlayerA.Draw()
	game.PlayCard(0, true)
}
