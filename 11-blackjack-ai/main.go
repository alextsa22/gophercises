package main

import (
	"fmt"
	"github.com/alextsa22/gophercises/11-blackjack-ai/blackjack"
)

func main() {
	game := blackjack.New()
	winnings := game.Play(blackjack.NewHumanAI())
	fmt.Printf("your balance: %d", winnings)
}
