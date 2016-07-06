package tajitsu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// NewGame creates and start a new game
func NewGame() {
	// Create decks of card for each player
	var deckA, deckB Deck

	// Add the combat cards
	f, e := ioutil.ReadFile("./combat_card.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	json.Unmarshal(f, &deckA)

	// Shuffle the deck
	deckA.Shuffle()
	deckB.Shuffle()
}
