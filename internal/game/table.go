package game

import (
	"log/slog"
	"net"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

type Table struct {
	ID      string
	Players []*player

	pc atomic.Int32
	mu sync.Mutex
}

func newTable() *Table {
	t := &Table{
		ID:      uuid.NewString(),
		Players: make([]*player, 0, 2), //  Heads-up only
	}

	slog.Info("new table created ", slog.String("TableID", t.ID))
	return t
}

func (t *Table) AddPlayer(conn net.Conn) {
	t.mu.Lock()
	defer t.mu.Unlock()

	newPlayer := newPlayer(conn)
	slog.Info("new player was added to table ", slog.String("TableID", t.ID), slog.String("PlayerID", newPlayer.ID))
	t.Players = append(t.Players, newPlayer)

	t.pc.Add(1)
}

func (t *Table) GetPlayersCount() int32 {
	return t.pc.Load()
}
