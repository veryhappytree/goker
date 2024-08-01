package game

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID    string
	Table *Table
}

func NewSession() *Session {
	s := &Session{
		ID:    uuid.NewString(),
		Table: newTable(),
	}

	slog.Info("new game session created", slog.String("SessionID", s.ID))
	return s
}

func (s *Session) AddPlayer(conn net.Conn) {
	s.Table.AddPlayer(conn)
}

func (s *Session) Run() {
	slog.Info("run game session", slog.String("SessionID", s.ID))

	for {
		select {
		case <-time.NewTicker(time.Duration(5) * time.Second).C:
			slog.Info("game session", slog.String("SessionID", s.ID), slog.Any("active players", s.Table.GetPlayersCount()))
		}
	}
}

func (s *Session) Close(wg *sync.WaitGroup) error {
	defer wg.Done()

	for _, p := range s.Table.Players {
		if err := p.disconnect(); err != nil {
			return err
		}
	}
	return nil
}
