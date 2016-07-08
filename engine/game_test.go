package engine

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var data []byte

func init() {
	// Load the combat cards
	f, e := ioutil.ReadFile("../data/combat_card.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	data = f
}

func TestPlayCard(t *testing.T) {
	game := NewGame(data)
	game.PlayerA.Draw()
	game.PlayCard(0, true, true)
}

func TestAbandonToVictory(t *testing.T) {
	game := NewGame(data)
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
	game := NewGame(data)
	if err := game.PlayCard(0, false, true); err != nil { // Player A plays a card
		t.Error(err)
	}
	if err := game.PlayCard(0, true, true); err != nil { // Player B plays a card on a new combo
		t.Error(err)
	}
	if err := game.Abandon(); err != nil { // Player A abandons
		t.Error(err)
	}
	if err := game.PickCombo(0); err != nil { // Player B picks combos
		t.Error(err)
	}
	if err := game.PickCombo(0); err != nil { // Player A picks combos
		t.Error(err)
	}
	if err := game.SetCurrentPlayer(true); err != nil {
		t.Error(err)
	}
	if game.State != Combat {
		t.Errorf("New game round has not began, game state is %s", game.State)
	}

	numberOfCardsPlayerA := len(game.PlayerA.Deck) + len(game.PlayerA.Hand)
	if numberOfCardsPlayerA != 9 {
		t.Errorf("Total number of cards in players deck and hand is %d, expected 9", numberOfCardsPlayerA)
	}
	numberOfCardsPlayerB := len(game.PlayerB.Deck) + len(game.PlayerB.Hand)
	if numberOfCardsPlayerB != 9 {
		t.Errorf("Total number of cards in players deck and hand is %d, expected 9", numberOfCardsPlayerB)
	}
}
