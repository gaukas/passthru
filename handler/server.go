package handler

import (
	"io"
	"net"

	"github.com/gaukas/passthru/config"
	"github.com/gaukas/passthru/protocol"
)

// A server, when started, will listen on a specific address specified in the config file
// Any incoming connections accepted on the listener will be COPIED
// and the copy will be passed to the protocol manager's FindAction method

type Server struct {
	serverAddr      config.ServerAddr
	listener        net.Listener
	protocolManager *protocol.ProtocolManager
}

// Required parameters will be provided from the main function
func NewServer(serverAddr config.ServerAddr, protocolManager *protocol.ProtocolManager) *Server {
	return &Server{
		serverAddr:      serverAddr,
		protocolManager: protocolManager,
	}
}

func (s *Server) Start() error {
	// create listener
	// start accepting connections （use goroutine）
	// for each connection, copy it and pass the copy to protocol manager

	cBuf := protocol.NewConnBuf()
	go io.Copy(cBuf, conn)

}
