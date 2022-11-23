package blackjack

import (
	"github.com/ethanefung/cards"
)

// Dealer is the Game Controller
type Dealer struct {
	deck  cards.Deck
	hand  Hand
	index int
	Game  *Game
}

// NewDealer returns a Dealer when given an instantiated game.
func NewDealer(g *Game) *Dealer {
	return &Dealer{Game: g}
}

// UseDecks will cause dealer to use a deck with 52 * n cards for gameplay.
func (d *Dealer) UseDecks(n int) {
	d.deck = cards.New()
	d.deck.Multiply(n)
	d.index = 0
}

// Shuffle will shuffle the deck the dealer is using.
func (d *Dealer) Shuffle(seed int64) {
	d.deck.Shuffle(seed)
}

// Deal will append the specified count of cards to all players within the Player's
// list. Deal will also append the specified count of card to the dealer's own hand.
func (d *Dealer) Deal(count int, p *PlayersList) {
	for i := 0; i < count; i++ {
		for curr := p; curr != nil; curr = curr.Tail {
			if d.index == len(d.deck) {
				return
			}
			curr.Head.Hand.Draw(d.deck[d.index])
			d.index++
		}
		d.hand.Draw(d.deck[d.index])
		d.index++
	}
}

// Hit appends one card from the dealers deck to the specified players hand if it is
// the player's turn
func (d *Dealer) Hit(p *Player) bool {
	if p != d.Game.Current.Head {
		return false
	}
	p.Hand.Draw(d.deck[d.index])
	d.index++
	return true
}

// Stay changes the game's current player to either the next player in the Players
// list or changes current to null.
func (d *Dealer) Stay() {
	d.Game.EndPlayerTurn()
}

// Play appends cards to the dealers hand as long as the value of the dealers hand is
// either below 17 or if the dealer has hand value of 17 and an ace.
func (d *Dealer) Play() {
	for d.hand.Value() < 17 || (d.hand.Value() == 17 && d.hand.HasAce()) {
		d.hand.Draw(d.deck[d.index])
		d.index++
	}
}

// Evaluate will change the game's current player if the current players hand value is
// over 21.
func (d *Dealer) Evaluate() {
	if d.Game.Current == nil {
		return
	}
	if d.Game.Current.Head.Hand.Value() > 21 {
		d.Game.Current = d.Game.Current.Tail
	}
}

// Clear removes all cards from players and dealer's hands.
func (d *Dealer) Clear() {
	curr := d.Game.Players
	for curr != nil {
		curr.Head.Hand = Hand{}
		curr = curr.Tail
	}

	d.hand = Hand{}
}

// ShowHand will return dealers full hand if all players have taken their turns for
// the round
func (d *Dealer) ShowHand() Hand {
	if d.Game.PlayersPlayed() {
		return d.hand
	}
	return d.hand[1:]
}
