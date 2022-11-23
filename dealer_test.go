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

	dealer.Hit(a)

	if len(a.Hand) != 3 {
		t.Fatalf("expected %v's hand to have 3 cards, but currently has %d", a.String(), len(a.Hand))

	}

	dealer.Stay()

	if game.Current.Head != b {
		t.Fatalf("expected in the given game state it is %v's turn but the current player is %v", b.String(), game.Current.Head.String())
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

	a.Hand = Hand{
		{Rank: cards.King},
		{Rank: cards.Nine},
		{Rank: cards.Nine},
	}

	dealer.Evaluate()

	if game.Current.Head == a {
		t.Fatalf("expected that the dealer would have moved on from player a but did not")
	}

	if game.Current.Head != b {
		t.Fatalf("expected that dealer waits on player two after player one busts")
	}

	game = New()

	dealer = game.Dealer

	a = NewPlayer("a")

	game.AddPlayer(a)

	game.Start()

	dealer.Evaluate()

	if game.Current.Head != a {
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