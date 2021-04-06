package blackjack

import (
	"fmt"

	deck "github.com/alextsa22/gophercises/09-deck"
)

type AI interface {
	Bet(shuffled bool) int
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hands [][]deck.Card, dealer []deck.Card)
}

type dealerAI struct{}

func (ai dealerAI) Bet(shuffled bool) int {
	return 1
}

func (ai dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dScore := Score(hand...)
	if dScore <= 16 || (dScore == 17 && Soft(hand...)) {
		return MoveHit
	}

	return MoveStand
}

func (ai dealerAI) Results(hands [][]deck.Card, dealer []deck.Card) {}

type humanAI struct{}

func NewHumanAI() AI {
	return humanAI{}
}

func (ai humanAI) Bet(shuffled bool) int {
	fmt.Println("what would you like to bet?")
	var bet int
	fmt.Scanf("%d\n", &bet)

	return bet
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	var input string
	for {
		fmt.Println("player:", hand)
		fmt.Println("dealer:", dealer)
		fmt.Println("what will you do? (h)it, (s)tand, (d)ouble, s(p)lit")

		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		case "d":
			return MoveDouble
		case "p":
			return MoveSplit
		default:
			fmt.Println("invalid option:", input)
			fmt.Println()
		}
	}
}

func (ai humanAI) Results(hands [][]deck.Card, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("player:")
	for _, h := range hands {
		fmt.Println("  ", h)
	}
	fmt.Println("dealer:", dealer)
	fmt.Println()
}
