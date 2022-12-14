package blackjack

import (
	"testing"

	"github.com/ethanefung/cards"
)

func TestGame(t *testing.T) {
	game := New()
	a, b := &Player{}, &Player{}
	game.AddPlayer(a)
	if game.Players == nil {
		t.Fatalf("No Players added")
	}
	if game.Players.Len() != 1 {
		t.Fatalf("Expected game to only one player but has %d", game.Players.Len())
	}
	curr := game.Players
	if curr.Head.Player != a {
		t.Fatalf("Expected player 'a' to be the first player, but was not")
	}
	game.AddPlayer(b)
	if game.Players.Len() != 2 {
		t.Fatalf("Expected game to have two players but has %d", game.Players.Len())
	}
	curr = curr.Tail
	if curr.Head.Player != b {
		t.Fatalf("Expected player 'b' to be the second player, but was not")
	}
}

func TestRemovePlayer(t *testing.T) {
	game := New()
	a, b, c := NewPlayer("a"), NewPlayer("b"), NewPlayer("c")

	game.AddPlayer(a)
	game.AddPlayer(b)
	game.AddPlayer(c)

	game.RemovePlayer(b)

	if game.Players.Len() != 2 {
		t.Fatalf("Attempted to remove one player of three, but was unsuccessful")
	}

	game.RemovePlayer(a)
	if game.Players.Len() != 1 {
		t.Fatalf("Attempted to remove the first player, but was unsuccessful")
	}

	game = New()

	game.AddPlayer(a)
	game.AddPlayer(a)
	game.AddPlayer(b)
	game.AddPlayer(a)
	game.RemovePlayer(a)
	if game.Players.Len() != 1 {
		t.Fatalf("Attempted to remove the first player added three times, but was unsuccessful")
	}

	game = New()

	game.AddPlayer(a)
	game.AddPlayer(a)
	game.AddPlayer(b)
	game.AddPlayer(a)
	game.AddPlayer(c)
	game.RemovePlayer(a)
	if game.Players.Len() != 2 {
		t.Fatalf("Attempted to remove player a, but was unsuccessful")
	}
}

func TestGameState(t *testing.T) {
	game := New()

	a, b := NewPlayer("a"), NewPlayer("b")
	dealer := game.Dealer

	game.AddPlayer(a)
	game.AddPlayer(b)

	val := game.Players

	val.Head.Hand = Hand{
		{Rank: cards.Ace},
		{Rank: cards.Jack},
	}

	val = val.Tail

	val.Head.Hand = Hand{
		{Rank: cards.Seven},
		{Rank: cards.Seven},
		{Rank: cards.Eight},
	}

	dealer.hand = Hand{
		{Rank: cards.Seven},
		{Rank: cards.King},
	}

	state := game.State()

	states := state.Players

	if states[a][0].Type != Win {
		t.Fatalf("expected that player 'a' to have won, but players state was %s", states[a][0].Type.String())
	}

	if states[b][0].Type != Bust {
		t.Fatalf("expected that player 'b' to have busted, but players state was %s", states[b][0].Type.String())
	}

	game.Start()
	state = game.State()

	states = state.Players

	if states[a][0].Type.String() != "Undetermined" {
		t.Fatalf("expected that player 'a's state to be Undetermined, but players state was %s", states[a][0].Type.String())
	}

	if states[b][0].Type.String() != "Undetermined" {
		t.Fatalf("expected that player 'b' to be Undetermined, but players state was %s", states[b][0].Type.String())
	}
}
