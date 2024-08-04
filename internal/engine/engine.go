package engine

import (
	"goker/internal/game"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
)

// dummy version
type engine struct {
	activeSessions map[string]*game.Session
	asc            atomic.Int32
	mu             sync.Mutex
	tcp            net.Listener
}

func Start() *engine {
	slog.Info("Start Goker Engine")

	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		slog.Error("init TCP listener was failed", slog.String("error", err.Error()))
		return nil
	}

	goker := engine{
		tcp:            listener,
		activeSessions: make(map[string]*game.Session),
	}

	go goker.seatPlayers()

	return &goker
}

func (g *engine) Shutdown() error {
	g.mu.Lock()

	wg := sync.WaitGroup{}
	wg.Add(len(g.activeSessions))

	for _, s := range g.activeSessions {
		go s.Close(&wg)
	}

	wg.Wait()

	g.mu.Unlock()

	err := g.tcp.Close()
	if err != nil {
		slog.Error("error shutting down engine", slog.String("error", err.Error()))
		return err
	}

	slog.Info("Shutdown engine")
	return nil
}

func (g *engine) seatPlayers() {
	for {
		conn, err := g.tcp.Accept()
		if err != nil {
			slog.Error("failed dial TCP connection", slog.String("error", err.Error()))
			continue
		}

		go g.play(conn)
	}
}

func (g *engine) play(conn net.Conn) {
	if g.asc.Load() == 0 {
		g.seatPlayerOnNewTable(conn)
		return
	}

	if table := g.findFreeTable(); table != nil {
		g.mu.Lock()
		table.AddPlayer(conn)
		g.mu.Unlock()
	} else {
		g.seatPlayerOnNewTable(conn)
	}
}

func (g *engine) seatPlayerOnNewTable(conn net.Conn) {
	g.mu.Lock()

	session := game.NewSession()
	session.AddPlayer(conn)

	g.activeSessions[session.ID] = session

	g.mu.Unlock()

	go session.Run()

	g.asc.Add(1)
}

func (g *engine) findFreeTable() *game.Table {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, session := range g.activeSessions {
		if session.Table.GetPlayersCount() == 1 {
			return session.Table
		}
	}
	return nil
}
