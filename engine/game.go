package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// NewGame creates and start a new game
func NewGame() {
	// Create the players
	var playerA, playerB Player

	// Add the combat cards
	f, e := ioutil.ReadFile("./data/combat_card.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var combatCards []*CombatCard
	json.Unmarshal(f, &combatCards)

	fmt.Printf("Combat cards: %s\n\n", combatCards)

	for _, combatCard := range combatCards {
		playerA.Deck = append(playerA.Deck, combatCard)
		playerB.Deck = append(playerB.Deck, combatCard)
	}

	// Shuffle the deck
	playerA.DeckShuffle()
	playerB.DeckShuffle()

}
