package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jglouis/tajitsu/engine"
)

var player engine.Player
var combatCards []*engine.CombatCard

func TestMain(m *testing.M) {
	f, e := ioutil.ReadFile("../data/combat_card.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	json.Unmarshal(f, &combatCards)
	for _, combatCard := range combatCards {
		player.Deck = append(player.Deck, combatCard)
	}

	os.Exit(m.Run())
}

func TestDraw(t *testing.T) {
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
