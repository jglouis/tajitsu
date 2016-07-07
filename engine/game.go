package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Game contains game related info (scores, cards in the arena)
type Game struct {
	State            GameState
	ScoreA, ScoreB   uint8 // score of each player, the first player winning two rounds, win the game
	PlayerA, PlayerB *Player
	CurrentPlayer    *Player          // The active player
	Combos           []CardCollection // the current combo is always the last one
	IsYinA           [][]bool         // Cards corresponding orientations, true id yin is pointing towards player A
}

// GameState gives indication on the game state
type GameState uint8

// GameState constants
//go:generate stringer -type=GameState
const (
	Combat  GameState = iota // Playing combat cards
	Recover                  // Picking the combos and adding the cards to the TestDeckShuffle
	ChoosingFirstPlayer
	End
)

func (game Game) getCurrentCombo() (CardCollection, []bool) {
	return game.Combos[len(game.Combos)-1], game.IsYinA[len(game.Combos)-1]
}

// PlayCard move a given card from current player's hand to the current combo
func (game *Game) PlayCard(pos uint8, isYinA bool) error {
	if game.State != Combat {
		return fmt.Errorf("Can only play a combat card, during the Combat step, current game state is %s", game.State)
	}

	player := game.CurrentPlayer
	card := player.Hand[pos]
	// Remove from hand
	copy(player.Hand[pos:], player.Hand[pos+1:])
	player.Hand[len(player.Hand)-1] = nil
	player.Hand = player.Hand[:len(player.Hand)-1]
	// Add to the current combo
	currentCombo, currentOrientations := game.getCurrentCombo()
	currentCombo = append(currentCombo, card.(*CombatCard))
	currentOrientations = append(currentOrientations, isYinA)

	return nil
}

// Abandon indicates the current player may or does not want to play a card
// The round is lost for him and the other player increments his score
func (game *Game) Abandon() error {
	if game.State != Combat {
		return fmt.Errorf("Can only abandon during the Combat step, current game state is %s", game.State)
	}

	// Increment score and check for victory
	if game.CurrentPlayer == game.PlayerA {
		game.ScoreB++
		if game.ScoreB == 2 {
			fmt.Println("Player B is victorious")
			game.State = End
			return nil
		}
	} else {
		game.ScoreA++
		if game.ScoreA == 2 {
			fmt.Println("Player A is victorious")
			game.State = End
			return nil
		}
	}

	// Switch game state to Recover or to ChoosingFirstPlayer if there are no combos to pick
	if len(game.Combos) == 0 || len(game.Combos[0]) == 0 {
		game.State = ChoosingFirstPlayer
	} else {
		game.State = Recover
	}

	return nil
}

// PickCombo moves all the card from the combo the current player's deck
func (game *Game) PickCombo(pos int) error {
	if game.State != Recover {
		return fmt.Errorf("Can only pick a combo, during the Recover step, current game state is %s", game.State)
	}
	if pos+1 > len(game.Combos) {
		return fmt.Errorf("The combo at position %d does not exist", pos)
	}
	// Add the cards of the combo to the current player's deck
	game.CurrentPlayer.Deck.Merge(game.Combos[pos])
	// Remove the combo from the slice
	game.Combos = append(game.Combos[:pos], game.Combos[pos+1:]...)

	// If all combos are picked, go to the next step
	if len(game.Combos) == 0 || len(game.Combos[0]) == 0 {
		game.State++
	}

	return nil
}

// SetCurrentPlayer allows to pick the new first player for the next round
// true for player A
func (game *Game) SetCurrentPlayer(isA bool) error {
	if game.State != ChoosingFirstPlayer {
		return fmt.Errorf("Can only set the first player, during the ChoosingFirstPlayer step, current game state is %s", game.State)
	}
	if isA {
		game.CurrentPlayer = game.PlayerA
	} else {
		game.CurrentPlayer = game.PlayerB
	}
	// Start the next round
	game.startNextRound()
	return nil
}
func (game *Game) startNextRound() {
	game.State = Combat
	game.IsYinA = [][]bool{[]bool{}}
	game.Combos = []CardCollection{CardCollection{}}
}

// NewGame creates and start a new game
func NewGame(dataPath string) *Game {
	game := new(Game)
	game.PlayerA = new(Player)
	game.PlayerB = new(Player)
	game.CurrentPlayer = game.PlayerA

	// Create the first combo
	game.startNextRound()

	// Load the combat cards
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
