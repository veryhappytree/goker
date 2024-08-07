package game

import (
	"log/slog"
	"net"

	"github.com/google/uuid"
)

type player struct {
	ID   string
	conn net.Conn
	// etc...
	Hand []*Card
}

func newPlayer(conn net.Conn) *player {
	p := &player{
		ID:   uuid.NewString(),
		conn: conn,
		Hand: make([]*Card, 0),
	}

	slog.Info("new player created ", slog.String("PlayerID", p.ID))
	return p
}

func (p *player) disconnect() error {
	if err := p.conn.Close(); err != nil {
		slog.Any("player disconnection error", err)
		return err
	}

	slog.Info("player disconnected", slog.String("PlayerID", p.ID))
	return nil
}
