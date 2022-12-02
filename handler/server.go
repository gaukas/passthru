package handler

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gaukas/passthru/config"
	"github.com/gaukas/passthru/internal/logger"
	"github.com/gaukas/passthru/protocol"
)

// A server, when started, will listen on a specific address specified in the config file
// Any incoming connections accepted on the listener will be COPIED
// and the copy will be passed to the protocol manager's FindAction method

type ServerMode uint8

const (
	SERVER_MODE_WORKER    ServerMode = iota // HandleNextConn() to be called by an external worker
	SERVER_MODE_UNLIMITED                   // unlimited number of connections will be handled
)

const (
	DEFAULT_TIMEOUT = 5 * time.Second
)

type Server struct {
	serverAddr config.ServerAddr
	listener   net.Listener

	protocolManager *protocol.ProtocolManager

	connBuf chan net.Conn
	mode    ServerMode
}

// Required parameters will be provided from the main function
func NewServer(serverAddr config.ServerAddr, protocolManager *protocol.ProtocolManager, mode ServerMode) *Server {
	return &Server{
		serverAddr:      serverAddr,
		protocolManager: protocolManager,
		mode:            mode,
	}
}

func (s *Server) Start() error {
	logger.Warnf("Starting server on %s", s.serverAddr)
	listener, err := net.Listen("tcp", s.serverAddr)
	if err != nil {
		logger.Errorf("Failed to start server on %s: %s", s.serverAddr, err)
		return err
	}

	s.listener = listener
	s.connBuf = make(chan net.Conn)

	go s.acceptLoop()

	return nil
}

func (s *Server) Stop() error {
	logger.Warnf("Stopping server on %s", s.serverAddr)
	close(s.connBuf)
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			return err
		}
	}
	logger.Infof("Server on %s stopped", s.serverAddr)
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			logger.Errorf("Failed to accept connection: %s, shutting down the server... ", err)
			s.Stop()
			return
		}
		logger.Infof("Accepted connection from %s", conn.RemoteAddr())

		if s.mode == SERVER_MODE_UNLIMITED {
			logger.Debugf("Starting a new goroutine to handle the connection from %s", conn.RemoteAddr())
			ctxExpire, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
			go s.handleConn(ctxExpire, conn)
			defer cancel()
		} else {
			logger.Debugf("Passing the connection from %s to the channel", conn.RemoteAddr())
			s.connBuf <- conn
		}
	}
}

// HandleNextConn() will block until a connection is available
// then call handleConn() to handle the connection upon it.
// Or it will return an error if the context is cancelled.
func (s *Server) HandleNextConn(ctx context.Context) error {
	select {
	case conn := <-s.connBuf:
		if conn == nil {
			logger.Errorf("Connection channel is closed, cannot handle the next connection")
			return ErrServerStopped
		}
		return s.handleConn(ctx, conn)
	case <-ctx.Done():
		logger.Errorf("Context is Done due to reason: %v, cannot handle the next connection", ctx.Err())
		return ctx.Err()
	}
}

func (s *Server) handleConn(ctx context.Context, conn net.Conn) error {
	defer conn.Close()
	wg := &sync.WaitGroup{}

	// Copy the connection
	// Pass the copy to the protocol manager
	// Get the action back
	// Perform the action
	cBuf := protocol.NewConnBuf()
	defer cBuf.Close()

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		io.Copy(cBuf, conn) // conn->cBuf
		conn.Close()
	}(wg)

	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithTimeout(ctx, DEFAULT_TIMEOUT)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()
	action, err := s.protocolManager.FindAction(ctx, cBuf)
	if err != nil && err != context.Canceled { // Canceled indicates a CATCHALL
		return err
	}

	switch action.Action {
	case config.ACTION_FORWARD:
		// dial up the destination
		connDst, err := net.Dial("tcp", action.ToAddr)
		if err != nil {
			return err
		}
		defer connDst.Close()

		logger.Infof("Forwarding connection from %s to %s", conn.RemoteAddr(), action.ToAddr)

		// Set downstream for the connection buffer
		err = cBuf.SetDownstream(connDst)
		if err != nil {
			return err
		}

		io.Copy(conn, connDst) // connDst->conn, so it is a bidirectional pipe
		wg.Wait()              // wait for conn->cBuf(->connDst) to finish
		return nil
	case config.ACTION_REJECT:
		return nil // do nothing, conn will be closed by defer
	default:
		return ErrUnknownAction
	}
}
