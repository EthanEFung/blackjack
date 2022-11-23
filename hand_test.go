package blackjack

import (
	"testing"

	"github.com/ethanefung/cards"
)

func TestHandValue(t *testing.T) {
	hand := Hand{
		{Rank: cards.Ace},
		{Rank: cards.King},
	}
	if hand.Value() != 21 {
		t.Fatalf("expected blackjack hand to have a value of 21, but got %d", hand.Value())
	}
	hand = Hand{
		{Rank: cards.Ace},
		{Rank: cards.King},
		{Rank: cards.Two},
	}
	if hand.Value() != 13 {
		t.Fatalf("expected hand to have a value of 13, but got %d", hand.Value())
	}
	hand = Hand{
		{Rank: cards.Ace},
		{Rank: cards.Ace},
		{Rank: cards.Ace},
		{Rank: cards.Two},
	}
	if hand.Value() != 15 {
		t.Fatalf("expected hand to have a value of 15, but got %d", hand.Value())
	}
	hand = Hand{
		{Rank: cards.Ace},
		{Rank: cards.King},
		{Rank: cards.Nine},
		{Rank: cards.Two},
	}
	if hand.Value() != 22 {
		t.Fatalf("expected hand to have a value of 22, but got %d", hand.Value())
	}
}
