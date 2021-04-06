package main

import (
	"fmt"
	"github.com/alextsa22/gophercises/11-blackjack-ai/blackjack"
)

func main() {
	opts := blackjack.Options{
		Decks:           3,
		Hands:           2,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(opts)
	winnings := game.Play(blackjack.NewHumanAI())
	fmt.Printf("your balance: %d", winnings)
}
