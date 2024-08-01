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

var goker engine

func Start() {
	slog.Info("Start Goker Engine")

	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		slog.Any("init TCP listener was failed", err)
		return
	}

	goker = engine{
		tcp:            listener,
		activeSessions: make(map[string]*game.Session),
		mu:             sync.Mutex{},
	}

	go goker.seatPlayers()
}

func Shutdown() error {
	goker.mu.Lock()

	wg := sync.WaitGroup{}
	wg.Add(len(goker.activeSessions))

	for _, s := range goker.activeSessions {
		go s.Close(&wg)
	}

	wg.Wait()

	goker.mu.Unlock()

	err := goker.tcp.Close()
	if err != nil {
		slog.Any("error shutting down engine", err)
		return err
	}

	slog.Info("Shutdown engine")
	return nil
}

func (g *engine) seatPlayers() {
	for {
		conn, err := goker.tcp.Accept()
		if err != nil {
			slog.Any("failed dial TCP connection", err)
			continue
		}

		go goker.play(conn)
	}
}

func (g *engine) play(conn net.Conn) {
	if g.asc.Load() == 0 {
		g.seatPlayerOnNewTable(conn)
		return
	}

	if table := g.findFreeTable(); table != nil {
		table.AddPlayer(conn)
	} else {
		g.seatPlayerOnNewTable(conn)
	}
}

func (g *engine) seatPlayerOnNewTable(conn net.Conn) {
	g.mu.Lock()

	session := game.NewSession()
	session.AddPlayer(conn)

	goker.activeSessions[session.ID] = session

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
