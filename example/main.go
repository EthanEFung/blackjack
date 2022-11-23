package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/ethanefung/blackjack"
)

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func clearTerminal() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}
}

func printState(state blackjack.GameState, g *blackjack.Game) {
	b := strings.Builder{}
	b.WriteString("Dealer: ")
	if g.PlayersPlayed() && state.Dealer.State != blackjack.Undetermined {
		b.WriteString(state.Dealer.State.String())
	}

	b.WriteRune('\n')

	if !g.PlayersPlayed() {
		b.WriteString("    --------------\n")
	}
	for _, card := range g.Dealer.ShowHand() {
		b.WriteString("    ")
		b.WriteString(card.String())
		b.WriteRune('\n')
	}
	b.WriteRune('\n')

	for curr := g.Players; curr != nil; curr = curr.Tail {
		p := curr.Head
		b.WriteString(p.String() + ": ")
		if g.PlayersPlayed() {
			b.WriteString(state.Players[p].State.String())
		}
		b.WriteRune('\n')

		for _, card := range p.Hand {
			b.WriteString("    ")
			b.WriteString(card.String())
			b.WriteRune('\n')
		}
		b.WriteRune('\n')
	}

	fmt.Print(b.String())
}

func main() {
	// give the user the ability to add players
	fmt.Print("Welcome to blackjack. To begin, type your name: ")

	reader := bufio.NewReader(os.Stdin)
	game := blackjack.New()
	dealer := game.Dealer

	for {
		text, _ := reader.ReadString('\n')
		input := strings.TrimSpace(text)
		if input == "" && game.Players == nil {
			fmt.Print("Any name would suffice: ")
			continue
		} else if input == "" {
			break
		}

		player := blackjack.NewPlayer(input)
		game.AddPlayer(player)
		fmt.Printf("Hi %s, shall we start? Type another name to add a player, or press enter to start...\n", input)
	}

	dealer.UseDecks(3)
	dealer.Shuffle(time.Now().Unix())
	fmt.Println("Let's begin.")

	for {
		dealer.Clear()
		dealer.Deal(2, game.Players)

		game.Start()

		for !game.PlayersPlayed() {
			clearTerminal()
			printState(game.State(), game)
			fmt.Print(game.Current.Head.Name, ", (h)it or (s)tay: ")
			text, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			option := strings.TrimSpace(text)
			if option == "h" {
				dealer.Hit(game.Current.Head)
			} else if option == "s" {
				dealer.Stay()
			}
			dealer.Evaluate()
		}

		dealer.Play()
		clearTerminal()
		printState(game.State(), game)
		fmt.Print("Again? Press enter to continue or type 'n + <enter>' to end: ")
		text, _ := reader.ReadString('\n')
		input := strings.TrimSpace(text)
		if input == "n" {
			break
		}
	}
	fmt.Println("\nThanks for playing.")
}
