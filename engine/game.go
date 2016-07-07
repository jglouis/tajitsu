package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Game contains game related info (scores, cards in the arena)
type Game struct {
	ScoreA, ScoreB   uint8 // score of each player, the first player winning two rounds, win the game
	PlayerA, PlayerB *Player
	CurrentPlayer    *Player         // The active player
	Combos           [][]*CombatCard // the current combo is always the last one
	IsYinUp          [][]bool        // Cards corresponding orientations
}

func (game Game) getCurrentCombo() []*CombatCard {
	return game.Combos[len(game.Combos)-1]
}

// PlayCard move a given card from current player's hand to the current combo
func (game *Game) PlayCard(pos uint8, isYinUp bool) {
	// Create a new combo if None exists
	if len(game.Combos) == 0 {
		game.Combos = [][]*CombatCard{[]*CombatCard{}}
	}

	player := game.CurrentPlayer
	card := player.Hand[pos]
	// Remove from hand
	copy(player.Hand[pos:], player.Hand[pos+1:])
	player.Hand[len(player.Hand)-1] = nil
	player.Hand = player.Hand[:len(player.Hand)-1]
	// Add to the current combo
	currentCombo := game.getCurrentCombo()
	currentCombo = append(currentCombo, card.(*CombatCard))
}

// NewGame creates and start a new game
func NewGame(dataPath string) *Game {
	game := new(Game)
	game.PlayerA = new(Player)
	game.PlayerB = new(Player)
	game.CurrentPlayer = game.PlayerA

	// Add the combat cards
	f, e := ioutil.ReadFile(dataPath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var combatCards []*CombatCard
	json.Unmarshal(f, &combatCards)

	fmt.Printf("Combat cards: %s\n\n", combatCards)

	for _, combatCard := range combatCards {
		game.PlayerA.Deck = append(game.PlayerA.Deck, combatCard)
		game.PlayerB.Deck = append(game.PlayerB.Deck, combatCard)
	}

	// Shuffle the deck
	game.PlayerA.DeckShuffle()
	game.PlayerB.DeckShuffle()

	for i := 0; i < 3; i++ {
		game.PlayerA.Draw()
		game.PlayerB.Draw()
	}

	fmt.Printf("Player A hand: %s\n", game.PlayerA.Hand)
	fmt.Printf("Player B hand: %s\n", game.PlayerB.Hand)

	return game
}
