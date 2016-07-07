package main

import "github.com/jglouis/tajitsu/engine"

func main() {
	game := engine.NewGame("./data/combat_card.json")

	game.PlayCard(0, false, true)
	game.Abandon()
	game.PickCombo(0)
	game.SetCurrentPlayer(true)

	game.PlayCard(0, false, true)
	game.Abandon()
	game.PickCombo(0)
	game.SetCurrentPlayer(true)
}
