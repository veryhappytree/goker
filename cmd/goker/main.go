package main

import (
	"fmt"
	"goker/internal/engine"
	"goker/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("♠♣goker♦♥")
	logger.Setup()

	go engine.Start()

	gracefulShutdown(engine.Shutdown)
}

func gracefulShutdown(ops ...func() error) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	if <-shutdown != nil {
		for _, op := range ops {
			if err := op(); err != nil {
				slog.Any("gracefullShutdown op failed", err)
			}
		}
	}
}
