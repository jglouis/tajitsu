package tajitsu

import "math/rand"

// Deck represents a deck of card
type Deck []Card

// DiscardPile represents a discard pile
// The bottome of the pile is index 0
type DiscardPile Deck

// Hand represents a hand of card visible to a player
type Hand Deck

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

// Discard put the nth card in the discard pile
func (hand Hand) Discard(n int, discardPile DiscardPile) {
	card := hand[n]
	copy(hand[n:], hand[n+1:])
	hand[len(hand)-1] = nil
	hand = hand[:len(hand)-1]
	discardPile = append(discardPile, card)
}
