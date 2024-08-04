package game

import (
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID    string
	Table *Table
	mu    sync.Mutex
	stop  chan struct{}
}

func NewSession() *Session {
	s := &Session{
		ID:    uuid.NewString(),
		Table: newTable(),
		stop:  make(chan struct{}),
	}

	slog.Info("new game session created", slog.String("SessionID", s.ID))
	return s
}

func (s *Session) AddPlayer(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Table.AddPlayer(conn)
}

func (s *Session) Run() {
	slog.Info("run game session", slog.String("SessionID", s.ID))
run:
	for {
		select {
		case <-time.NewTicker(time.Duration(1) * time.Second).C:
			slog.Info("game session", slog.String("SessionID", s.ID), slog.Any("active players", s.Table.GetPlayersCount()))
		case <-s.stop:
			break run
		}
	}
}

func (s *Session) Close(wg *sync.WaitGroup) error {
	defer wg.Done()

	close(s.stop)

	for _, p := range s.Table.Players {
		if err := p.disconnect(); err != nil {
			return err
		}
	}
	return nil
}
