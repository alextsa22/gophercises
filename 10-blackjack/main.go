package main

import (
	"fmt"
	"strings"

	deck "github.com/alextsa22/gophercises/09-deck"
)

func main() {
	var gs GameState
	gs = Shuffle(gs)

	for i := 0; i < 10; i++ {
		gs = Deal(gs)

		var input string
		for gs.State == StatePlayerTurn {
			fmt.Println("player:", gs.Player)
			fmt.Println("dealer:", gs.Dealer.DealerString())
			fmt.Println("what will you do? (h)it, (s)tand")
			fmt.Scanf("%s\n", &input)
			switch input {
			case "h":
				gs = Hit(gs)
			case "s":
				gs = Stand(gs)
			default:
				fmt.Println("invalid option:", input)
				fmt.Println()
			}
		}

		for gs.State == StateDealerTurn {
			if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
				gs = Hit(gs)
			} else {
				gs = Stand(gs)
			}
		}

		gs = EndHand(gs)
	}
}

type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}

	return strings.Join(strs, ", ")
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}

	for _, c := range h {
		if c.Rank == deck.Ace {
			return minScore + 10
		}
	}

	return minScore
}

func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}

	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func Shuffle(gs GameState) GameState {
	tmp := clone(gs)
	tmp.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return tmp
}

func Deal(gs GameState) GameState {
	tmp := clone(gs)
	tmp.Player = make(Hand, 0, 5)
	tmp.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, tmp.Deck = draw(tmp.Deck)
		tmp.Player = append(tmp.Player, card)
		card, tmp.Deck = draw(tmp.Deck)
		tmp.Dealer = append(tmp.Dealer, card)
	}

	tmp.State = StatePlayerTurn
	return tmp
}

func Stand(gs GameState) GameState {
	tmp := clone(gs)
	tmp.State++
	return tmp
}

func Hit(gs GameState) GameState {
	tmp := clone(gs)
	hand := tmp.CurrentPlayer()

	var card deck.Card
	card, tmp.Deck = draw(tmp.Deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		return Stand(tmp)
	}

	return tmp
}

func EndHand(gs GameState) GameState {
	tmp := clone(gs)
	pScore, dScore := tmp.Player.Score(), tmp.Dealer.Score()

	fmt.Println("==FINAL HANDS==")
	fmt.Println("player:", tmp.Player, "\nscore:", pScore)
	fmt.Println("dealer:", tmp.Dealer, "\nscore:", dScore)

	switch {
	case pScore > 21:
		fmt.Println("you busted")
	case dScore > 21:
		fmt.Println("dealer busted")
	case pScore > dScore:
		fmt.Println("you win!")
	case dScore > pScore:
		fmt.Println("you lose")
	case dScore == pScore:
		fmt.Println("draw")
	}
	fmt.Println()

	tmp.Player = nil
	tmp.Dealer = nil
	return tmp
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

type State int8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

func clone(gs GameState) GameState {
	tmp := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}

	copy(tmp.Deck, gs.Deck)
	copy(tmp.Player, gs.Player)
	copy(tmp.Dealer, gs.Dealer)

	return tmp
}
