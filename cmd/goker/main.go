package main

import (
	"fmt"
	"goker/internal/deck"
	"goker/pkg/logger"
)

func main() {
	fmt.Println("goker")

	logger.Setup()

	d := deck.New()

	d.Print()

	d.Shuffle()

	d.Print()
}
