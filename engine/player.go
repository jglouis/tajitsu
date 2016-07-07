package engine

import "math/rand"

// Player represents a player and contains all things belonging to one player
type Player struct {
	Hand, Deck, DiscardPile CardCollection
}

// CardCollection is a slice of Card
type CardCollection []Card

// DeckShuffle randomizes the player's deck
func (player *Player) DeckShuffle() {
	for i := range player.Deck {
		j := rand.Intn(i + 1)
		player.Deck[i], player.Deck[j] = player.Deck[j], player.Deck[i]
	}
}

// Draw pop the first card from the deck to the player's hand
// Returns the drawn card or nil if the deck was empty
func (player *Player) Draw() Card {
	if len(player.Deck) == 0 {
		return nil
	}
	var card Card
	card, player.Deck = player.Deck[len(player.Deck)-1], player.Deck[:len(player.Deck)-1]
	player.Hand = append(player.Hand, card)
	return card
}

// Discard put the n_th card from player's hand to his discard pile
// if the card does not exists (or if the hand is empty), does nothing
func (player *Player) Discard(n int) {
	if n+1 > len(player.Hand) {
		return
	}
	card := player.Hand[n]
	copy(player.Hand[n:], player.Hand[n+1:])
	player.Hand[len(player.Hand)-1] = nil
	player.Hand = player.Hand[:len(player.Hand)-1]
	player.DiscardPile = append(player.DiscardPile, card)
}
