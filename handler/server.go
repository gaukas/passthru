package handler

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gaukas/passthru/config"
	"github.com/gaukas/passthru/protocol"
)

// A server, when started, will listen on a specific address specified in the config file
// Any incoming connections accepted on the listener will be COPIED
// and the copy will be passed to the protocol manager's FindAction method

type Server struct {
	serverAddr config.ServerAddr
	listener   net.Listener

	timeout         time.Duration
	protocolManager *protocol.ProtocolManager
}

// Required parameters will be provided from the main function
func NewServer(serverAddr config.ServerAddr, protocolManager *protocol.ProtocolManager) *Server {
	return &Server{
		serverAddr:      serverAddr,
		protocolManager: protocolManager,
		timeout:         5 * time.Second,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.serverAddr)
	if err != nil {
		return err
	}

	s.listener = listener

	go s.acceptLoop()

	return nil
}

func (s *Server) Stop() error {
	return s.listener.Close()
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.listener.Close()
			return
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Copy the connection
	// Pass the copy to the protocol manager
	// Get the action back
	// Perform the action
	cBuf := protocol.NewConnBuf()
	defer cBuf.Close()

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		io.Copy(cBuf, conn)
		cBuf.Close()
		conn.Close()
	}(wg)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	action, err := s.protocolManager.FindAction(ctxTimeout, cBuf)
	if err != nil {
		return
	}

	switch action.Action {
	case config.ACTION_FORWARD:
		// dial up the destination
		connDst, err := net.Dial("tcp", action.ToAddr)
		if err != nil {
			return
		}
		defer connDst.Close()

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			io.Copy(connDst, cBuf) // anything read from cBuf will be written to connDst
		}(wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			io.Copy(conn, connDst) // anything read from connDst will be written to conn
		}(wg)

		wg.Wait() // wait for all goroutines to finish

		return
	case config.ACTION_REJECT:
		return // do nothing, conn will be closed by defer
	default:
		return
	}
}
