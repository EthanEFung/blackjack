/*

A simple library of reasonable and extendable utilities to build blackjack applications.

*/
package blackjack

// WinType is the state used to determine if a player or dealer has
// won in the round.
type WinType int

const (
    // Undetermined means the round is not yet over.
    Undetermined WinType = iota
    // Lose means that the value of the hand is inferior.
    Lose 
    // Win means that the value of the hand is superior.
    Win 
    // Bust is a losing WinType that means the value of the hand is over 21.
    Bust
    // Push means the value of the hand is the same as the dealers.
    Push
)
//go:generate stringer -type=WinType

// WinState is the state of a player or dealer and the value of the hand.
type WinState struct {
    // Value is the numeric value of the hand the state represents.
    Value int
    // State is the WinType of the player or dealer.
    State WinType
}

// GameState is the aggregate of players' and dealer's winstates.
type GameState struct {
    // Dealer is the WinState of the game's dealer.
    Dealer WinState
    // Players is a map of Player WinStates.
    Players map[*Player]WinState
}

// PlayersList is a link list of players.
type PlayersList struct {
    // Head will always represent a Player struct.
	Head *Player
    // Tail will always represent another PlayerList. 
	Tail *PlayersList
}

// Len returns the length of the players list.
func (p *PlayersList) Len() int {
	if p.Tail == nil {
		return 1
	}
	return 1 + p.Tail.Len()
}

// Game is aggregate struct of PlayersList and Dealer.
type Game struct {
    // Players is the PlayersList that represents all the players who are current in game.
	Players *PlayersList
    // Current is the PlayersList that represents the player who is currently playing their turn.
	Current *PlayersList
    // Dealer is the controller of the Game.
	Dealer  *Dealer
}

// New returns an instance of Game.
func New() *Game {
	game := new(Game)
	dealer := NewDealer(game)
	game.Dealer = dealer
	return game
}

// AddPlayer appends the specified player to the end of Game.Players.
func (g *Game) AddPlayer(p *Player) {
	if g.Players == nil {
		g.Players = &PlayersList{
			Head: p,
		}
	} else {
		// TODO: optimize
		curr := g.Players
		for curr.Tail != nil {
			curr = curr.Tail
		}
		curr.Tail = &PlayersList{
			Head: p,
		}
	}
}

// RemovePlayer removes the player from the Game.Players list.
func (g *Game) RemovePlayer(p *Player) {
	curr := g.Players
	for curr.Head == p {
		curr = curr.Tail
		g.Players = curr
	}

	for curr.Tail != nil {
		if curr.Tail.Head == p && curr.Tail.Tail == nil {
			curr.Tail = nil
			break
		} else if curr.Tail.Head == p {
			curr.Tail = curr.Tail.Tail
		}
		curr = curr.Tail
	}
}

// Start assigns the current Game.Players list to Game.Current (necessary for dealer to
// know which player to deal to).
func (g *Game) Start() bool {
	if g.Players == nil {
		return false
	}
	g.Current = g.Players
	return true
}

// EndPlayerTurn will change Game.Current to the next player within the Game.Players
// list. May also assign Game.Current to nil if there are no more players.
func (g *Game) EndPlayerTurn() bool {
	if g.Current.Tail == nil {
		g.Current = nil
		return false
	}
	g.Current = g.Current.Tail
	return true
}

// State returns a GameState struct which changes depending on whether the game is done.
func (g *Game) State() GameState {
    done := g.PlayersPlayed()
    dealer := g.Dealer

    dealerState := WinState{}
    // evaluate dealer's hand
    if dealer.hand.Value() > 21 {
        dealerState.State = Bust
    } else if done {
        dealerState.Value = dealer.hand.Value()
    }

    winStates := make(map[*Player]WinState, g.Players.Len())

    for list := g.Players; list != nil; list = list.Tail {
        player := list.Head
        winState := WinState{ Value: player.Hand.Value() }
        playerHand, dealerHand := player.Hand.Value(), dealer.hand.Value()
        if !done {
            winStates[player] = winState
            continue
        }
        // game is done
        if  playerHand > 21 {
            winState.State = Bust
        } else if playerHand == dealerHand {
            winState.State = Push
        } else if playerHand > dealerHand || dealerHand > 21 {
            winState.State = Win
        } else {
            winState.State = Lose
        }

        winStates[player] = winState
    }

    return GameState{
        Dealer: dealerState,
        Players: winStates,
    }
}

// PlayersPlayed returns false if there are still players who have yet to take their
// turn for the round.
func (g *Game) PlayersPlayed() bool {
    return g.Current == nil
}
