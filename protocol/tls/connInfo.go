package tls

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"time"

	"github.com/gaukas/passthru/protocol"
	tls "github.com/refraction-networking/utls"
)

type ConnInfo struct {
	SNI  string
	ALPN string
}

// Check https://github.com/refraction-networking/utls/blob/2179f286686bdd60b90151993024fb9cfc21420b/conn.go#L991
func ParseClientHello(ctx context.Context, cbuf *protocol.ConnBuf) (ConnInfo, error) {
	for ctx.Err() == nil {
		// peek first 5 bytes
		buf := make([]byte, 5)
		err := cbuf.Peek(buf, 5)
		if err != nil {
			if err == io.EOF {
				return ConnInfo{}, err
			} else {
				time.Sleep(20 * time.Millisecond)
				continue
			}
		}

		if buf[0] != 0x16 {
			return ConnInfo{}, errors.New("not a TLS connection")
		}
		// parse length
		length := binary.BigEndian.Uint16(buf[3:])
		// peek the rest of the message
		buf = make([]byte, length+5)
		err = cbuf.Peek(buf, int(length+5))
		if err != nil {
			if err == protocol.ErrNotEnoughData {
				continue
			} else {
				return ConnInfo{}, err
			}
		}

		// check if it's a client hello
		if buf[5] != 0x01 {
			return ConnInfo{}, errors.New("not start with a client hello")
		}

		clientHello := tls.UnmarshalClientHello(buf[5:]) // drop the first 5 bytes as TLS record header
		if clientHello == nil {
			return ConnInfo{}, errors.New("failed to parse client hello")
		}

		ci := ConnInfo{
			SNI: clientHello.ServerName,
		}
		if len(clientHello.AlpnProtocols) > 0 {
			ci.ALPN = clientHello.AlpnProtocols[0]
		}

		return ci, nil
	}

	return ConnInfo{}, ctx.Err()
}
