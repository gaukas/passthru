package tls_test

import (
	"context"
	"testing"

	"github.com/gaukas/passthru/protocol"
	"github.com/gaukas/passthru/protocol/tls"
)

var (
	sampleClientHello = []byte{
		0x16, 0x03, 0x01, 0x01, 0xb6, 0x01, 0x00, 0x01,
		0xb2, 0x03, 0x03, 0x0e, 0x96, 0x92, 0x12, 0x90,
		0xa1, 0x32, 0x1e, 0x3d, 0xbd, 0x05, 0xc1, 0xb7,
		0x20, 0xae, 0x69, 0x59, 0x40, 0xc0, 0x68, 0xc9,
		0xa7, 0x08, 0xa0, 0x10, 0x1a, 0xf1, 0x28, 0xa0,
		0x60, 0xf6, 0x8a, 0x20, 0xee, 0x83, 0xcf, 0x77,
		0x62, 0xbf, 0xb9, 0xef, 0xa4, 0x8f, 0x34, 0x59,
		0x8e, 0x9f, 0x1a, 0x30, 0x05, 0x51, 0x30, 0xce,
		0xa6, 0x54, 0xae, 0xbb, 0x1f, 0xdb, 0xd0, 0xe3,
		0xc0, 0x72, 0xdd, 0x52, 0x00, 0x8c, 0x1a, 0x1a,
		0xc0, 0x12, 0xc0, 0x13, 0xc0, 0x07, 0xc0, 0x27,
		0xcc, 0x14, 0xc0, 0x2f, 0x13, 0x01, 0xc0, 0x14,
		0x13, 0x02, 0xc0, 0x28, 0xcc, 0xa9, 0xc0, 0x30,
		0xc0, 0x73, 0xc0, 0x60, 0xc0, 0x72, 0xc0, 0x61,
		0xc0, 0x2c, 0xc0, 0x76, 0xc0, 0xaf, 0xc0, 0x77,
		0xc0, 0xad, 0xcc, 0xa8, 0xc0, 0x24, 0x13, 0x05,
		0xc0, 0x0a, 0x13, 0x04, 0xc0, 0x2b, 0x13, 0x03,
		0xc0, 0xae, 0xcc, 0x13, 0xc0, 0xac, 0xc0, 0x11,
		0xc0, 0x23, 0x00, 0x0a, 0xc0, 0x09, 0x00, 0x2f,
		0xc0, 0x08, 0x00, 0x3c, 0x00, 0x9a, 0xc0, 0x9c,
		0x00, 0xc4, 0xc0, 0xa0, 0x00, 0x88, 0x00, 0x9c,
		0x00, 0xbe, 0x00, 0x35, 0x00, 0x45, 0x00, 0x3d,
		0x00, 0x9f, 0xc0, 0x9d, 0xc0, 0xa3, 0xc0, 0xa1,
		0xc0, 0x9f, 0x00, 0x9d, 0x00, 0x6b, 0x00, 0x41,
		0x00, 0x39, 0x00, 0xba, 0x00, 0x9e, 0x00, 0x84,
		0xc0, 0xa2, 0x00, 0xc0, 0xc0, 0x9e, 0x00, 0x07,
		0x00, 0x67, 0x00, 0x04, 0x00, 0x33, 0x00, 0x05,
		0x00, 0x16, 0x01, 0x00, 0x00, 0xdd, 0x3a, 0x3a,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x11,
		0x00, 0x00, 0x0e, 0x31, 0x30, 0x37, 0x2e, 0x31,
		0x38, 0x32, 0x2e, 0x32, 0x36, 0x2e, 0x31, 0x31,
		0x39, 0x00, 0x17, 0x00, 0x00, 0x00, 0x01, 0x00,
		0x01, 0x01, 0xff, 0x01, 0x00, 0x01, 0x00, 0x00,
		0x0a, 0x00, 0x0a, 0x00, 0x08, 0x00, 0x1d, 0x00,
		0x17, 0x00, 0x18, 0x00, 0x19, 0x00, 0x0b, 0x00,
		0x02, 0x01, 0x00, 0x00, 0x23, 0x00, 0x00, 0x00,
		0x10, 0x00, 0x3c, 0x00, 0x3a, 0x02, 0x68, 0x71,
		0x03, 0x68, 0x32, 0x63, 0x02, 0x68, 0x32, 0x06,
		0x73, 0x70, 0x64, 0x79, 0x2f, 0x33, 0x06, 0x73,
		0x70, 0x64, 0x79, 0x2f, 0x32, 0x06, 0x73, 0x70,
		0x64, 0x79, 0x2f, 0x31, 0x08, 0x68, 0x74, 0x74,
		0x70, 0x2f, 0x31, 0x2e, 0x31, 0x08, 0x68, 0x74,
		0x74, 0x70, 0x2f, 0x31, 0x2e, 0x30, 0x08, 0x68,
		0x74, 0x74, 0x70, 0x2f, 0x30, 0x2e, 0x39, 0x00,
		0x0d, 0x00, 0x14, 0x00, 0x12, 0x04, 0x03, 0x08,
		0x04, 0x04, 0x01, 0x05, 0x03, 0x08, 0x05, 0x05,
		0x01, 0x08, 0x06, 0x06, 0x01, 0x02, 0x01, 0x00,
		0x33, 0x00, 0x2b, 0x00, 0x29, 0xba, 0xba, 0x00,
		0x01, 0x00, 0x00, 0x1d, 0x00, 0x20, 0x2d, 0xe3,
		0x2c, 0x09, 0xf3, 0x3b, 0x19, 0xbc, 0xd3, 0x1e,
		0x5a, 0x3a, 0x14, 0xba, 0x1a, 0xd6, 0x24, 0x45,
		0x7d, 0x89, 0x28, 0xee, 0x51, 0x75, 0x1a, 0x4b,
		0x56, 0xc1, 0xa8, 0x7f, 0xb4, 0x59, 0x00, 0x2d,
		0x00, 0x02, 0x01, 0x01, 0x00, 0x2b, 0x00, 0x0b,
		0x0a, 0x1a, 0x1a, 0x03, 0x04, 0x03, 0x03, 0x03,
		0x02, 0x03, 0x01,
	}
)

func TestParseClientHello(t *testing.T) {
	ctx := context.Background()
	cBuf := protocol.NewConnBuf()
	cBuf.Write(sampleClientHello)

	connInfo, err := tls.ParseClientHello(ctx, cBuf)
	if err != nil {
		t.Fatalf("ParseClientHello failed: %v", err)
	}
	if connInfo.SNI != "107.182.26.119" {
		t.Fatalf("SNI mismatch: %v", connInfo.SNI)
	}

	if connInfo.ALPN != "hq" {
		t.Fatalf("ALPN mismatch: %v", connInfo.ALPN)
	}
}
