package main

import (
	"fmt"
	"github.com/alextsa22/gophercises/11-blackjack-ai/blackjack"

	deck "github.com/alextsa22/gophercises/09-deck"
)

type basicAI struct {
}

func (ai *basicAI) Bet(shuffled bool) int {
	return 100
}

func (ai *basicAI) Play(hand []deck.Card, dealer deck.Card) blackjack.Move {
	score := blackjack.Score(hand...)
	if len(hand) == 2 {
		if hand[0] == hand[1] {
			cardScore := blackjack.Score(hand[0])
			if cardScore >= 8 && cardScore != 10 {
				return blackjack.MoveSplit
			}
		}

		if (score == 10 || score == 11) && !blackjack.Soft(hand...) {
			return blackjack.MoveDouble
		}
	}

	dealerScore := blackjack.Score(dealer)
	if dealerScore >= 5 && dealerScore <= 6 {
		return blackjack.MoveStand
	}

	if score < 13 {
		return blackjack.MoveHit
	}

	return blackjack.MoveStand
}

func (ai *basicAI) Results(hands [][]deck.Card, dealer []deck.Card) {}

func main() {
	opts := blackjack.Options{
		Decks:           4,
		Hands:           1000,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(opts)
	winnings := game.Play(new(basicAI))
	fmt.Printf("your balance: %d", winnings)
}
