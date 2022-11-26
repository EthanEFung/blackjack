package blackjack

import (
	"github.com/ethanefung/cards"
)

// Dealer is the Game Controller
type Dealer struct {
	deck     cards.Deck
	hand     Hand
	index    int
	Game     *Game
	Splitter PlayersList
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

// Hit will add a card for to the current players hand.
func (d *Dealer) Hit() bool {
	if d.Game.Current == nil {
		return false
	}
	listVal := d.Game.Current.Head
    listVal.Hand.Draw(d.deck[d.index])
	d.index++
	return true
}

// Stay changes the game's current player to either the next player in the Players
// list or changes current to null.
func (d *Dealer) Stay() {
	d.Game.EndPlayerTurn()
}

// Surrender will first subtract half of the wager amount from the current players
// winnings and set the bet amount to zero. Surrender will also end the players turn.
func (d *Dealer) Surrender() bool {
	if d.Game.Current == nil {
		return false
	}
	player := d.Game.Current.Head
	half := player.Wager / 2
	d.Game.Current.Head.Player.Winnings -= half
	d.Game.Current.Head.Wager = 0
	d.Game.EndPlayerTurn()
	return true
}

// Double will multiply the current players wager by 2 and hit if the player has only
// two cards.
// ! To note that this shouldn't be allowed on split hands
func (d *Dealer) Double() bool {
	if d.Game.Current == nil {
		return false
	}
	player := d.Game.Current.Head
	if len(player.Hand) > 2 {
		return false
	}
	d.Game.Current.Head.Wager *= 2
	d.Hit()
	d.Game.EndPlayerTurn()
	return true
}

// Bet will change the Wager of the player to the specified amount.
func (d *Dealer) Bet(listVal *ListVal, wager int) bool {
	if listVal == nil {
		return false
	}
	listVal.Wager = wager
	return true
}

// Collect resolves all game Players Winnings based on the state of the game.
// TODO: splitting has ramifications for how collection works. Where the hands
// that were split also have to be taken into account.
func (d *Dealer) Collect() {
	states := d.Game.State().Players
	for curr := d.Game.Players; curr != nil; curr = curr.Tail {
        listVal := curr.Head
        for _, state := range states[listVal.Player] {
            switch state.Type {
            case Push:
                continue
            case Win:
                listVal.Player.Winnings += curr.Head.Wager
            case Lose:
                listVal.Player.Winnings -= curr.Head.Wager
            case Bust:
                listVal.Player.Winnings -= curr.Head.Wager
            default:
                continue
            }

        }
	}
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
// TODO: Evaluate also might need to be reimagined as we have to make sure that the player
// has no more hands to be evaluated before moving on to the next player
func (d *Dealer) Evaluate() {
	if d.Game.Current == nil {
		return
	}
	if d.Game.Current.Head.Hand.Value() > 21 {
		d.Game.Current = d.Game.Current.Tail
	}
}

// Clear removes all cards from players and dealer's hands.
// TODO: how does split effect clearing of cards?
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
