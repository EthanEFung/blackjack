package blackjack

import (
	"fmt"
)

// Player is a struct representing an end user in the game.
type Player struct {
	// Name is the name of the player.
	Name string
	// Winnings is the current number of bets the user has currently
	Winnings int
}

// NewPlayer returns an reference to a player with the specified name.
func NewPlayer(name string) *Player {
	return &Player{Name: name}
}

func (p *Player) String() string {
	return fmt.Sprintf("Player(%s)", p.Name)
}
