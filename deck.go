package tajitsu

import "math/rand"

// Deck represents a deck of card
type Deck []Card

// Shuffle randomizes the deck
func (deck Deck) Shuffle() {
	for i := range deck {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}

// Draw pop the first card from the deck
func (deck Deck) Draw() Card {
	card, deck := deck[len(deck)-1], deck[:len(deck)-1]
	return card
}
