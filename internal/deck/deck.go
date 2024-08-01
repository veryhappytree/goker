package deck

import (
	"fmt"
	"time"

	"math/rand"
)

type cardSuit string

const (
	Clubs    cardSuit = "clubs"
	Diamonds cardSuit = "diamonds"
	Hearts   cardSuit = "hearts"
	Spades   cardSuit = "spades"
)

func (c cardSuit) prettyCardSuit() string {
	switch c {
	case Clubs:
		return "♣"
	case Diamonds:
		return "♦"
	case Hearts:
		return "♥"
	case Spades:
		return "♠"
	default:
		panic("?")
	}
}

type cardValue int

const (
	Two cardValue = iota + 1
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func (c cardValue) prettyCardValue() string {
	switch c {
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "10"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		panic("?")
	}
}

var cardSuites = [4]cardSuit{Clubs, Diamonds, Hearts, Spades}

const totalCards = 52

type Card struct {
	suit  cardSuit
	value cardValue
}

type Deck struct {
	Cards [totalCards]Card
}

func New() *Deck {
	var cards [totalCards]Card

	i := 0
	for _, s := range cardSuites {
		for v := Two; v <= Ace; v++ {
			c := Card{suit: s, value: v}
			cards[i] = c
			i++
		}
	}

	return &Deck{Cards: cards}
}

func (d Deck) Print() {
	fmt.Printf("Total cards - %d\n", len(d.Cards))

	i := 0
	for _, card := range d.Cards {
		i++
		fmt.Printf("%s%s", card.value.prettyCardValue(), card.suit.prettyCardSuit())
		if i == 13 {
			fmt.Println("")
			i = 0
		}
	}
}

func (d *Deck) Shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for i := len(d.Cards) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}
