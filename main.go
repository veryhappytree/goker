package main

import (
	"fmt"
	"goker/deck"
)

func main() {
	fmt.Println("goker")

	d := deck.New()

	d.Print()

	d.Shuffle()

	d.Print()
}
