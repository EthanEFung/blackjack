package blackjack

import (
	"github.com/ethanefung/cards"
)

var cardValues = map[cards.Rank]int{
	cards.Two:   2,
	cards.Three: 3,
	cards.Four:  4,
	cards.Five:  5,
	cards.Six:   6,
	cards.Seven: 7,
	cards.Eight: 8,
	cards.Nine:  9,
	cards.Ten:   10,
	cards.Jack:  10,
	cards.Queen: 10,
	cards.King:  10,
	cards.Ace:   11,
}

// Hand is a slice of cards.
type Hand []cards.Card

// Draw will append the specified card to the hand.
func (h *Hand) Draw(c cards.Card) {
	*h = append(*h, c)
}

// Value returns the integer value of the hand.
func (h *Hand) Value() int {
	var nAces int
	var val int

	for _, card := range *h {
		if card.Rank == cards.Ace {
			nAces++
		}
		val += cardValues[card.Rank]
	}

	for nAces > 0 && val > 21 {
		val -= 10
		nAces--
	}
	return val
}

// HasAce will return true if the hand has an ace.
func (h *Hand) HasAce() bool {
	for _, card := range *h {
		if card.Rank == cards.Ace {
			return true
		}
	}
	return false
}
