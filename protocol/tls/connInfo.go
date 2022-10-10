package tls

import "github.com/gaukas/passthru/protocol"

type ConnInfo struct {
	SNI  string
	ALPN string
}

func ParseConnInfo(cBuf *protocol.ConnBuf) (ConnInfo, error) {
	// TODO: identify SNI and ALPN
	return ConnInfo{}, nil
}
