package blackjack

import (
	deck "github.com/alextsa22/gophercises/09-deck"
)

type AI interface {
	Bet() int
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hand [][]deck.Card, dealer []deck.Card)
}
