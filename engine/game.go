package engine

import (
	"encoding/json"
	"errors"
	"fmt"
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

// PlayCard move a given card from current player's hand to the current combo
func (game *Game) PlayCard(pos uint8, newCombo bool, isYinA bool) error {
	if game.State != Combat {
		return fmt.Errorf("Can only play a combat card, during the Combat step, current game state is %s", game.State)
	}

	// Start a new combo if necessary
	if newCombo && len(game.Combos[len(game.Combos)-1]) != 0 {
		game.Combos = append(game.Combos, CardCollection{})
		game.IsYinA = append(game.IsYinA, []bool{})
	}

	player := game.CurrentPlayer
	card := player.Hand[pos].(*CombatCard)

	lastComboIndex := len(game.Combos) - 1

	// Check if it is a valid play
	if !newCombo && len(game.Combos[lastComboIndex]) != 0 {
		// 1) Check card orientation
		lastIsYinA := game.IsYinA[lastComboIndex][len(game.IsYinA[lastComboIndex])-1]
		if isYinA != lastIsYinA {
			return errors.New("Wrong card orientation")
		}

		// 2) Check opponent stance
		lastCard := game.Combos[lastComboIndex][len(game.Combos[lastComboIndex])-1]
		var lastOpponentStance Stance
		var futureOpponentStance Stance
		if lastIsYinA && game.CurrentPlayer == game.PlayerA {
			lastOpponentStance = lastCard.(*CombatCard).YangStance
		} else {
			lastOpponentStance = lastCard.(*CombatCard).YinStance
		}
		if isYinA && game.CurrentPlayer == game.PlayerA {
			futureOpponentStance = card.YangStance
		} else {
			futureOpponentStance = card.YinStance
		}
		if futureOpponentStance != lastOpponentStance {
			return fmt.Errorf("Opponent stance not respected: is %s, but would be %s with the card was played", lastOpponentStance, futureOpponentStance)
		}
	}
	// Remove from hand
	copy(player.Hand[pos:], player.Hand[pos+1:])
	player.Hand[len(player.Hand)-1] = nil
	player.Hand = player.Hand[:len(player.Hand)-1]

	// Add to the current combo
	game.Combos[lastComboIndex] = append(game.Combos[lastComboIndex], card)
	game.IsYinA[lastComboIndex] = append(game.IsYinA[lastComboIndex], isYinA)

	// Resolve effects
	var playerYin, playerYang *Player
	if isYinA {
		playerYin = game.PlayerA
		playerYang = game.PlayerB
	} else {
		playerYin = game.PlayerB
		playerYang = game.PlayerA
	}
	for n, effect := range []Effect{card.YinEffect, card.YangEffect} {
		var affectedPlayer *Player
		switch Balance(n + 1) {
		case Yin:
			affectedPlayer = playerYin
		case Yang:
			affectedPlayer = playerYang
		}
		switch effect {
		case Draw:
			affectedPlayer.Draw()
		case Discard:
			affectedPlayer.Discard(0) // should ask what to Discard
			// TODO handle aspects
		}
	}

	game.switchCurrentPlayer()

	return nil
}

func (game *Game) switchCurrentPlayer() {
	if game.CurrentPlayer == game.PlayerA {
		game.CurrentPlayer = game.PlayerB
	} else {
		game.CurrentPlayer = game.PlayerA
	}
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
		game.switchCurrentPlayer()
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
	game.CurrentPlayer.Deck = append(game.CurrentPlayer.Deck, game.Combos[pos]...)
	// Remove the combo from the slice
	game.Combos = append(game.Combos[:pos], game.Combos[pos+1:]...)

	// If all combos are picked, go to the next step
	if len(game.Combos) == 0 || len(game.Combos[0]) == 0 {
		game.State++
	}

	game.switchCurrentPlayer()

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
	// All cards in deck and in hand goes to their player's respective deck
	for _, player := range []*Player{game.PlayerA, game.PlayerB} {
		player.Deck = append(player.Deck, player.Hand...)
		player.Hand = CardCollection{}
		player.Deck = append(player.Deck, player.DiscardPile...)
		player.DiscardPile = CardCollection{}
	}

	// Reset game state and combo area
	game.State = Combat
	game.IsYinA = [][]bool{[]bool{}}
	game.Combos = []CardCollection{CardCollection{}}

	// Shuffle the deck
	game.PlayerA.DeckShuffle()
	game.PlayerB.DeckShuffle()

	// Draw five cards
	for i := 0; i < 5; i++ {
		game.PlayerA.Draw()
		game.PlayerB.Draw()
	}
}

// NewGame creates and start a new game
func NewGame(data []byte) *Game {
	game := new(Game)
	game.PlayerA = new(Player)
	game.PlayerB = new(Player)
	game.CurrentPlayer = game.PlayerA

	var combatCards []*CombatCard
	json.Unmarshal(data, &combatCards)

	for _, combatCard := range combatCards {
		game.PlayerA.Deck = append(game.PlayerA.Deck, combatCard)
		game.PlayerB.Deck = append(game.PlayerB.Deck, combatCard)
	}

	game.startNextRound()

	fmt.Printf("Player A hand: %s\n", game.PlayerA.Hand)
	fmt.Printf("Player B hand: %s\n", game.PlayerB.Hand)

	return game
}
