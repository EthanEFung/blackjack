package blackjack

import (
	"testing"

	"github.com/ethanefung/cards"
)

func TestNewDealerSetup(t *testing.T) {
	game := New()
	a, b, c := &Player{}, &Player{}, &Player{}
	game.AddPlayer(a)
	game.AddPlayer(b)
	game.AddPlayer(c)

	dealer := game.Dealer

	dealer.UseDecks(1)
	dealer.Shuffle(0)
	dealer.Deal(2, game.Players)

	curr := game.Players
	for curr != nil {
		if len(curr.Head.Hand) != 2 {
			t.Fatalf("attempted to deal every player 2 cards but got %d, %v", len(curr.Head.Hand), curr.Head.Hand)
		}
		curr = curr.Tail
	}

	curr = game.Players
	kinds := make(map[cards.Card]bool)

	for curr != nil {
		p := curr.Head
		for _, card := range p.Hand {
			kinds[card] = true
		}
		curr = curr.Tail
	}

	for _, card := range dealer.hand {
		kinds[card] = true
	}

	if len(dealer.hand) != 2 {
		t.Fatalf("expected dealer's hand to have 2 cards but has %d cards", len(dealer.hand))
	}

	if len(kinds) != 8 {
		t.Fatalf("expected 8 unique cards to be dealt to 3 players and a dealer but got %d", len(kinds))
	}
}

func TestNewDealerGamePlay(t *testing.T) {
	game := New()
	a := &Player{Name: "a"}
	b := &Player{Name: "b"}
	c := &Player{Name: "c"}

	game.AddPlayer(a)
	game.AddPlayer(b)
	game.AddPlayer(c)
	dealer := game.Dealer

	dealer.UseDecks(1)
	dealer.Shuffle(0)
	dealer.Deal(2, game.Players)

	game.Start()

	dealer.Hit()

	if game.Current.Head.Player != a {
		t.Fatalf("expected the current player to be player a, but it was not")
	}

	firstHand := game.Current.Head.Hand

	if len(firstHand) != 3 {
		t.Fatalf("expected %v's hand to have 3 cards, but currently has %d", a.String(), len(firstHand))
	}

	dealer.Stay()

	if game.Current.Head.Player != b {
		t.Fatalf("expected in the given game state it is %v's turn but the current player is %v", b.String(), game.Current.Head.Player.String())
	}
}

func TestDealerPlay(t *testing.T) {
	game := New()

	dealer := game.Dealer
	dealer.deck = cards.Deck{
		{Rank: cards.Jack},
		{Rank: cards.Ace},
	}

	dealer.Play()

	if dealer.hand.Value() != 21 {
		t.Fatalf("expected that given only kings and aces, the dealer would have drawn either 20 or 21, value is %d", dealer.hand.Value())
	}

	game = New()
	dealer = game.Dealer

	dealer.deck = cards.Deck{
		{Rank: cards.Six},
		{Rank: cards.Ace},
		{Rank: cards.Ten},
		{Rank: cards.Four},
	}

	dealer.Play()

	if dealer.hand.Value() != 6+1+10+4 {
		t.Fatalf("expected the dealer would keep drawing on a soft 17, but something went wrong, %d", dealer.hand.Value())
	}
}

func TestDealerClear(t *testing.T) {
	game := New()

	dealer := game.Dealer

	a, b, c := NewPlayer("a"), NewPlayer("b"), NewPlayer("c")
	game.AddPlayer(a)
	game.AddPlayer(b)
	game.AddPlayer(c)

	dealer.UseDecks(1)

	dealer.Deal(2, game.Players)

	for curr := game.Players; curr != nil; curr = curr.Tail {
		player := curr.Head
		if len(player.Hand) != 2 {
			t.Fatalf("expected that dealer would give each player 2 cards but instead dealt %d", len(player.Hand))
		}
	}

	dealer.Clear()

	for curr := game.Players; curr != nil; curr = curr.Tail {
		player := curr.Head
		if len(player.Hand) != 0 {
			t.Fatalf("expected dealer to have removed cards from players hand but hand has %d cards", len(player.Hand))
		}
	}
}

func TestDealerEvaluate(t *testing.T) {
	game := New()

	dealer := game.Dealer

	a, b := NewPlayer("a"), NewPlayer("b")

	game.AddPlayer(a)
	game.AddPlayer(b)

	game.Start() // assigns current

	game.Current.Head.Hand = Hand{
		{Rank: cards.King},
		{Rank: cards.Nine},
		{Rank: cards.Nine},
	}

	dealer.Evaluate()

	if game.Current.Head.Player == a {
		t.Fatalf("expected that the dealer would have moved on from player a but did not")
	}

	if game.Current.Head.Player != b {
		t.Fatalf("expected that dealer waits on player two after player one busts")
	}

	game = New()

	dealer = game.Dealer

	a = NewPlayer("a")

	game.AddPlayer(a)

	game.Start()

	dealer.Evaluate()

	if game.Current.Head.Player != a {
		t.Fatalf("expected that since players hand is less than 21, that the current player is still a")
	}
	dealer.Stay()

	dealer.Evaluate()
	if game.Current != nil {
		t.Fatalf("expected that since the player chose to stay, there would be no current player")
	}
}

func TestDealerShowHand(t *testing.T) {
	game := New()

	game.AddPlayer(NewPlayer("a"))
	game.Dealer.hand = Hand{
		{Rank: cards.Ace},
		{Rank: cards.Jack},
	}

	game.Start()
	shown := game.Dealer.ShowHand()

	if len(shown) != 1 {
		t.Fatalf("expected that dealer to only show one card if players have not played but shows %d", len(shown))
	}

	game.Dealer.Stay()
	shown = game.Dealer.ShowHand()

	if len(shown) != 2 {
		t.Fatalf("expected that dealer to only show two card if all players have played but shows %d", len(shown))
	}
}

func TestDealerBet(t *testing.T) {
	game := New()

	a := NewPlayer("a")
	game.AddPlayer(a)
	game.Start()
	game.Dealer.Bet(game.Current.Head, 2)

	if game.Current.Head.Wager != 2 {
		t.Fatalf("expected the dealer to set players a's wager to 2 but wager is %d", game.Current.Head.Wager)
	}
}

func TestDealerSurrender(t *testing.T) {
	game := New()

	a := NewPlayer("a")
	game.AddPlayer(a)
	game.Start()

	game.Dealer.Bet(game.Current.Head, 2)

	game.Start()
	game.Dealer.Surrender()

	if a.Winnings != -1 {
		t.Fatalf("expected half the wager of player a to be subtracted from winnings, but player's winnings were %d", a.Winnings)
	}

	if game.Players.Head.Wager != 0 {
		t.Fatalf("expected after surrender that player a's wager would be zero but wager was %d", game.Current.Head.Wager)
	}
}

func TestDealerDouble(t *testing.T) {
	game := New()
	a := NewPlayer("a")
	game.AddPlayer(a)
	game.Dealer.UseDecks(1)
	game.Start()

	currentVal := game.Current.Head

	game.Dealer.Bet(currentVal, 2)
	game.Dealer.Deal(2, game.Players)
	game.Start()

	game.Dealer.Double()

	if currentVal.Wager != 4 {
		t.Fatalf("expected the players wager to have doubled but got a wager of %d", game.Current.Head.Wager)
	}

	if len(currentVal.Hand) != 3 {
		t.Fatalf("expected having doubled, the dealer to have dealt another card to player but player has %d cards", len(currentVal.Hand))
	}

	if game.Current != nil {
		t.Fatalf("expected that users turn would have ended after doubling, but it did not")
	}
}

func TestDealerCollect(t *testing.T) {
	game := New()
	a, b, c := NewPlayer("a"), NewPlayer("b"), NewPlayer("c")
	game.AddPlayer(a)
	game.AddPlayer(b)
	game.AddPlayer(c)

	for val := game.Players; val != nil; val = val.Tail {
		game.Dealer.Bet(val.Head, 2)
	}

	val := game.Players

	val.Head.Hand = Hand{
		{Rank: cards.Ace},
		{Rank: cards.Jack},
	}

	val = val.Tail

	val.Head.Hand = Hand{
		{Rank: cards.Jack},
		{Rank: cards.Six},
	}

	val = val.Tail

	val.Head.Hand = Hand{
		{Rank: cards.Jack},
		{Rank: cards.Seven},
	}

	game.Dealer.hand = Hand{
		{Rank: cards.Jack},
		{Rank: cards.Seven},
	}

	game.Dealer.Collect()

	if a.Winnings != 2 {
		t.Fatalf("expected player a who has won to have 2 winnings, but has %d", a.Winnings)
	}
	if b.Winnings != -2 {
		t.Fatalf("expected player b who has lost to have -2 winnings, but has %d", b.Winnings)
	}
	if c.Winnings != 0 {
		t.Fatalf("expected player c who has pushed to have 0 winnings, but has %d", b.Winnings)
	}
}

func TestDealerSplitAndCollect(t *testing.T) {
	game := New()
	a := NewPlayer("a")
	game.AddPlayer(a)

	game.Dealer.UseDecks(1)
	game.Dealer.Shuffle(2)

	game.Start()

	game.Dealer.Bet(game.Players.Head, 2)

	aceOfSpades := cards.Card{Rank: cards.Ace, Suit: cards.Spades}
	aceOfHearts := cards.Card{Rank: cards.Ace, Suit: cards.Hearts}

	game.Players.Head.Hand = Hand{aceOfSpades, aceOfHearts}

	game.Dealer.Split()

	if game.Current.Head.Hand.Value() != 19 {
		t.Fatalf("something happened with the current test setup, expected the current players hand to have a value of 13 but instead was: %d", game.Current.Head.Hand.Value())
	}
	game.Dealer.Stay()

	if game.Current.Head.Hand.Value() != 13 {
		t.Fatalf("something happened with the current test setup, expected the current players hand to have a value of 19 but instead was: %d", game.Current.Head.Hand.Value())
	}

	game.Dealer.Stay()

	game.Dealer.Play()

	game.Dealer.Collect()

	dealerHand := game.Dealer.ShowHand()

	if dealerHand.Value() != 20 {
		t.Fatalf("something happened with the current test setup, expected the dealers hand to have a value of 20 but instead was: %d", dealerHand.Value())
	}

	if game.Players.Head.Player.Winnings != -4 {
		t.Fatalf("expected players hands to result in -4 winnings (lost both) but instead was %d", game.Players.Head.Player.Winnings)
	}
}

func TestDealerSplit(t *testing.T) {
	game := New()
	a, b := NewPlayer("a"), NewPlayer("b")
	game.AddPlayer(a)
	game.AddPlayer(b)

	game.Dealer.UseDecks(1)
	game.Dealer.Shuffle(0)

	game.Start()

	aceOfSpades := cards.Card{Rank: cards.Ace, Suit: cards.Spades}
	aceOfHearts := cards.Card{Rank: cards.Ace, Suit: cards.Hearts}

	game.Players.Head.Hand = Hand{
		aceOfSpades,
		aceOfHearts,
	}

	game.Dealer.Split()

	if len(game.Current.Head.Hand) != 2 {
		t.Fatalf("expected after split that the current players hand be dealt one card to the first hand but currently has %d", len(game.Current.Head.Hand))
	}

	for _, card := range game.Current.Head.Hand {
		if card == aceOfHearts {
			t.Fatalf("expected that the current players hand to be split, but the hand contained a card it should not have")
		}
	}

	game.Dealer.Stay()

	if len(game.Current.Head.Hand) != 2 {
		t.Fatalf("expected after split that the current players hand be dealt one card to their second hand but currently has %d", len(game.Current.Head.Hand))
	}

	var hasAceOfHearts bool
	for _, card := range game.Current.Head.Hand {
		if card == aceOfHearts {
			hasAceOfHearts = true
		}
	}

	if !hasAceOfHearts {
		t.Fatalf("expected hand to have an ace of hearts because the hand was split by the player before")
	}

	if game.Current.Head.Player != a {
		t.Fatalf("expected the hand to be associated with player a but was not")
	}

	game.Dealer.Stay()

	if game.Current.Head.Player != b {
		t.Fatalf("expectdd that the split still keeps reference to the previous tail, but did not")
	}
}

func TestDealerNoSplit(t *testing.T) {
	game := New()
	a, b := NewPlayer("a"), NewPlayer("b")
	game.AddPlayer(a)
	game.AddPlayer(b)

	game.Dealer.UseDecks(1)
	game.Dealer.Shuffle(0)

	game.Start()

	game.Dealer.Deal(2, game.Players)

	game.Dealer.Split()

	if game.Players.Len() > 2 {
		t.Fatalf("expected that on a hand without a pair, split would be unsuccessful")
	}
}

func TestDealerResetTable(t *testing.T) {
	game := New()
	a, b := NewPlayer("a"), NewPlayer("b")
	game.AddPlayer(a)
	game.AddPlayer(b)

	game.Dealer.UseDecks(1)
	game.Dealer.Shuffle(0)

	eightOfSpades := cards.Card{Rank: cards.Eight, Suit: cards.Spades}
	eightOfClubs := cards.Card{Rank: cards.Eight, Suit: cards.Clubs}
	aceOfSpades := cards.Card{Rank: cards.Ace, Suit: cards.Spades}

	game.Start()

	firstPlayerVal := game.Current.Head
	secondPlayerVal := game.Current.Tail.Head

	game.Current.Head.Hand = Hand{
		eightOfSpades,
		eightOfClubs,
	}

	game.Current.Tail.Head.Hand = Hand{
		eightOfClubs,
		aceOfSpades,
	}

	game.Dealer.Split()

	if game.Current.Tail.Head == secondPlayerVal {
		t.Fatalf("expected that a dealer splitting would result in the insertion of a new list val but did not happen")
	}

	game.Dealer.Clear()
	game.Dealer.ResetTable()

	if game.Current.Head != firstPlayerVal {
		t.Fatalf("expected that the first players default list val to remain but the reference was lost")
	}

	if game.Current.Tail.Head != secondPlayerVal {
		t.Fatalf("expected that after the dealer.ResetTable is called that the split list val would be removed but it was not")
	}
}
