package main

import "github.com/jglouis/tajitsu/engine"

func main() {

	data, err := Asset("data/combat_card.json")
	if err != nil {
		panic(err)
	}

	game := engine.NewGame(data)

	game.PlayCard(0, false, true)
	game.Abandon()
	game.PickCombo(0)
	game.SetCurrentPlayer(true)

	game.PlayCard(0, false, true)
	game.Abandon()
	game.PickCombo(0)
	game.SetCurrentPlayer(true)
}
