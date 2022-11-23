package blackjack

import (
	"testing"
)

func TestPlayer(t *testing.T) {
	p := NewPlayer("John Doe")

	if p.Name != "John Doe" {
		t.Fatalf("created a player with name 'John Doe', but name was %s", p.Name)
	}

	if p.String() != "Player(John Doe)" {
		t.Fatalf("unexpected player.String() result %v", p.String())
	}

}
