package tajitsu

// CombatCard represents a classic combat card
type CombatCard struct {
	YinStance, YangStance                                          Stance
	Advantage                                                      Balance
	YinEffect, YangEffect, YinAdvantageEffect, YangAdvantageEffect Effect
}

// Balance represents the orientaion of a card
type Balance int

// Balance constants, Yin or Yang
const (
	Yin Balance = iota
	Yang
)

// Stance represents the combattant stance
type Stance int

// Stance constants
const (
	Serpent Stance = iota
	Tiger
	Dragon
)

// Effect represents a card effect
type Effect int

// Effect constants
const (
	Draw Effect = iota
	Discard
	SwitchAspect
)
