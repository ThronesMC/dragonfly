package dragonfly

import (
	"errors"
	"fmt"
	"github.com/dragonfly-tech/dragonfly/dragonfly/player"
	"github.com/dragonfly-tech/dragonfly/dragonfly/session"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sirupsen/logrus"
	"log"
)

// Server implements a Dragonfly server. It runs the main server loop and handles the connections of players
// trying to join the server.
type Server struct {
	c        Config
	log      *logrus.Logger
	listener *minecraft.Listener
	players  chan *player.Player
}

// New returns a new server using the Config passed. If nil is passed, a default configuration is returned.
// (A call to dragonfly.DefaultConfig().)
// The Logger passed will be used to log errors and information to. If nil is passed, a default Logger is
// used by calling logrus.New().
func New(c *Config, log *logrus.Logger) *Server {
	if log == nil {
		log = logrus.New()
	}
	players := make(chan *player.Player)
	if c == nil {
		return &Server{c: DefaultConfig(), log: log, players: players}
	}
	return &Server{c: *c, log: log, players: players}
}

// Accept accepts an incoming player into the server. It blocks until a player connects to the server.
// Accept returns an error if the Server is closed using a call to Close.
func (server *Server) Accept() (*player.Player, error) {
	p, ok := <-server.players
	if !ok {
		return nil, errors.New("server closed")
	}
	return p, nil
}

// Run runs the server and blocks until it is closed using a call to Close(). When called, the server will
// accept incoming connections.
// After a call to Run, calls to Server.Accept() may be made to accept players into the server.
func (server *Server) Run() error {
	if err := server.startListening(); err != nil {
		return err
	}
	server.run()
	return nil
}

// Start runs the server but does not block, unlike Run, but instead accepts connections on a different
// goroutine. Connections will be accepted until the listener is closed using a call to Close.
// One started, players may be accepted using Server.Accept().
func (server *Server) Start() error {
	if err := server.startListening(); err != nil {
		return err
	}
	go server.run()
	return nil
}

// startListening starts making the Minecraft listener listen, accepting new connections from players.
func (server *Server) startListening() error {
	server.log.Info("Starting server...")

	w := server.log.Writer()
	defer func() {
		_ = w.Close()
	}()

	server.listener = &minecraft.Listener{
		// We wrap a log.Logger around our Logrus logger so that it will print in the same format as the
		// normal Logrus logger would.
		ErrorLog:       log.New(w, "", 0),
		ServerName:     server.c.Server.Name,
		MaximumPlayers: server.c.Server.MaximumPlayers,
	}
	if err := server.listener.Listen("raknet", server.c.Network.Address); err != nil {
		return fmt.Errorf("listening on address failed: %v", err)
	}

	server.log.Infof("Server running on %v.\n", server.c.Network.Address)
	return nil
}

// run runs the server, continuously accepting new connections from players. It returns when the server is
// closed by a call to Close.
func (server *Server) run() {
	for {
		c, err := server.listener.Accept()
		if err != nil {
			// Accept will only return an error if the Listener was closed, meaning trying to continue
			// listening is futile.
			return
		}
		go server.handleConn(c.(*minecraft.Conn))
	}
}

// handleConn handles an incoming connection accepted from the Listener.
func (server *Server) handleConn(conn *minecraft.Conn) {
	data := minecraft.GameData{WorldName: server.c.Server.WorldName}
	if err := conn.StartGame(data); err != nil {
		return
	}
	p := &player.Player{}
	s := session.New(p, conn, server.log)
	*p = *player.NewWithSession(s)
	s.Handle()

	server.players <- p
}

// Close closes the server, making any call to Run/Accept cancel immediately.
func (server *Server) Close() error {
	close(server.players)
	return server.listener.Close()
}