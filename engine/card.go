package engine

import (
	"encoding/json"
	"fmt"
)

// Card is a generic card
type Card interface {
	fmt.Stringer
}

// CombatCard represents a classic combat card
type CombatCard struct {
	YinStance, YangStance                 Stance
	Advantage                             Balance
	YinEffect, YangEffect                 Effect
	YinFinisherEffect, YangFinisherEffect []Effect
	json.Unmarshaler
}

func (card CombatCard) String() string {
	return fmt.Sprintf("Combat card %s (Yin) / %s (Yang)", card.YinStance, card.YangStance)
}

func stringToStance(str string) Stance {
	switch str {
	case "Snake":
		return Snake
	case "Tiger":
		return Tiger
	case "Dragon":
		return Dragon
	}
	return -1
}

func stringToEffect(str string) Effect {
	switch str {
	case "Draw":
		return Draw
	case "Discard":
		return Discard
	case "SwitchAspect":
		return SwitchAspect
	}
	return NoEffect
}

func (card *CombatCard) UnmarshalJSON(data []byte) error {
	var v struct {
		YinStance, YangStance, Advantage, YinEffect, YangEffect string
		YinFinisherEffect, YangFinisherEffect                   []string
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	card.YinStance = stringToStance(v.YinStance)
	card.YangStance = stringToStance(v.YangStance)
	card.Advantage = Equi
	card.YinEffect = stringToEffect(v.YinEffect)
	card.YangEffect = stringToEffect(v.YangEffect)
	card.YinFinisherEffect = []Effect{}
	for _, str := range v.YinFinisherEffect {
		card.YinFinisherEffect = append(card.YinFinisherEffect, stringToEffect(str))
	}
	card.YangFinisherEffect = []Effect{}
	for _, str := range v.YangFinisherEffect {
		card.YangFinisherEffect = append(card.YangFinisherEffect, stringToEffect(str))
	}

	return nil
}

// Balance represents the orientaion of a card
type Balance uint8

// Balance constants, Yin or Yang
const (
	Equi Balance = iota
	Yin
	Yang
)

// Stance represents the combattant stance
type Stance int8

// Stance constants
//go:generate stringer -type=Stance
const (
	Snake Stance = iota
	Tiger
	Dragon
)

// Effect represents a card effect
type Effect uint8

// Effect constants
const (
	NoEffect Effect = iota
	Draw
	Discard
	SwitchAspect
)
