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
	json.Unmarshal(f, &(playerA.Deck))
	json.Unmarshal(f, &(playerA.Deck))

	fmt.Println(playerA.Deck)

	// Shuffle the deck
	playerA.DeckShuffle()
	playerB.DeckShuffle()

}
