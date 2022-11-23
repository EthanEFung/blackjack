package blackjack

import (
	"fmt"
)

// Player is a struct representing an end user in the game.
type Player struct {
    // Name is the name of the player.
	Name string
    // Hand is the slice of cards the user owns in the given round.
	Hand Hand
}

// NewPlayer returns an reference to a player with the specified name.
func NewPlayer(name string) *Player {
	return &Player{Name: name}
}

func (p *Player) String() string {
	return fmt.Sprintf("Player(%s)", p.Name)
}
