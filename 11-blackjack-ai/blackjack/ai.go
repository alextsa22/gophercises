package blackjack

import (
	"fmt"

	deck "github.com/alextsa22/gophercises/09-deck"
)

type AI interface {
	Bet() int
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hand [][]deck.Card, dealer []deck.Card)
}

type dealerAI struct{}

func (ai dealerAI) Bet() int {
	return 1
}

func (ai dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dScore := Score(hand...)
	if dScore <= 16 || (dScore == 17 && Soft(hand...)) {
		return MoveHit
	}

	return MoveStand
}

func (ai dealerAI) Results(hand [][]deck.Card, dealer []deck.Card) {}

type humanAI struct{}

func NewHumanAI() AI {
	return humanAI{}
}

func (ai humanAI) Bet() int {
	return 1
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	var input string
	for {
		fmt.Println("player:", hand)
		fmt.Println("dealer:", dealer)
		fmt.Println("what will you do? (h)it, (s)tand")

		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		default:
			fmt.Println("invalid option:", input)
			fmt.Println()
		}
	}
}

func (ai humanAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("player:", hand)
	fmt.Println("dealer:", dealer)
	fmt.Println()
}
