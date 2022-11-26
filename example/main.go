package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
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
	if g.PlayersPlayed() && state.Dealer.Type != blackjack.Undetermined {
		b.WriteString(state.Dealer.Type.String())
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
		val := curr.Head
		b.WriteString(val.Player.String() + ": ")
		if g.PlayersPlayed() {
			b.WriteString(state.Players[val.Player][0].Type.String())
		}
		b.WriteRune('\n')

		for _, card := range val.Hand {
			b.WriteString("    ")
			b.WriteString(card.String())
			b.WriteRune('\n')
		}
		b.WriteRune('\n')
	}

	fmt.Print(b.String())
}

func readStdin(reader *bufio.Reader) (string, error) {
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input := strings.TrimSpace(text)
	return input, nil
}

func main() {
	// give the user the ability to add players
	fmt.Print("Welcome to blackjack. To begin, type your name: ")

	reader := bufio.NewReader(os.Stdin)
	game := blackjack.New()
	dealer := game.Dealer

	for {
		input, err := readStdin(reader)
		if err != nil {
			fmt.Println("Couldn't read your input")
			continue
		}
		if input == "" && game.Players == nil {
			fmt.Print("Any name would suffice: ")
			continue
		} else if input == "" {
			break
		}

		player := blackjack.NewPlayer(input)
		game.AddPlayer(player)
		fmt.Printf("Hi %s, shall we start? Type another name to add a player, or press enter to start: ", input)
	}

	dealer.UseDecks(3)
	dealer.Shuffle(time.Now().Unix())
	fmt.Println("Let's begin.")

	for {
		fmt.Printf("First, everyone place bets\n")

		for curr := game.Players; curr != nil; curr = curr.Tail {
			val := curr.Head

			fmt.Printf("%s current winnings: %d\nwager: ", val.Player.String(), val.Player.Winnings)
			for {
				input, err := readStdin(reader)
				if err != nil {
					fmt.Printf("trouble reading your wager: ")
					continue
				}
				bet, err := strconv.Atoi(input)
				if err != nil {
					fmt.Printf("please enter a valid integer: ")
					continue
				}
				dealer.Bet(val, bet)
				break
			}
		}

		dealer.Clear()
		dealer.Deal(2, game.Players)

		game.Start()

		for !game.PlayersPlayed() {
			clearTerminal()
			printState(game.State(), game)
			fmt.Print(game.Current.Head.Player.Name, ", (h)it or (s)tay: ")
			option, err := readStdin(reader)
			if err != nil {
				fmt.Println("couldn't read your input")
			}
			if option == "h" {
				dealer.Hit()
			} else if option == "s" {
				dealer.Stay()
			}
			dealer.Evaluate()
		}

		dealer.Play()
		dealer.Collect()
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
