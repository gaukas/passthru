package tls_test

import (
	"context"
	"testing"
	"time"

	"github.com/gaukas/passthru/config"
	"github.com/gaukas/passthru/protocol"
	"github.com/gaukas/passthru/protocol/tls"
)

var (
	tlsProtocol = tls.Protocol{}
)

func TestProtocol(t *testing.T) {
	testApplyRules(t)
	testIdentify(t)
	testIdentifyExpired(t)
	testIdentifyWithNotEnoughData(t)
}

func testApplyRules(t *testing.T) {
	rules := []config.Rule{
		"CATCHALL",
		"SNI cloudflare-dns.com",
		"SNI dns.quad9.net",
		"ALPN h2",
	}

	err := tlsProtocol.ApplyRules(rules)
	if err != nil {
		t.Errorf("Error applying rules: %s", err)
	}
}

func testIdentify(t *testing.T) {
	var cBuf *protocol.ConnBuf

	ctx := context.Background()

	cBuf = protocol.NewConnBuf()
	cBuf.Write(CH_cloudflare_dns_com)
	rule, err := tlsProtocol.Identify(ctx, cBuf)
	if err != nil {
		t.Errorf("Error identifying rule: %s", err)
	}
	if rule != "SNI cloudflare-dns.com" {
		t.Errorf("Wrong rule identified: %s", rule)
	}

	cBuf = protocol.NewConnBuf()
	cBuf.Write(CH_quad9)
	rule, err = tlsProtocol.Identify(ctx, cBuf)
	if err != nil {
		t.Errorf("Error identifying rule: %s", err)
	}
	if rule != "SNI dns.quad9.net" {
		t.Errorf("Wrong rule identified: %s", rule)
	}

	cBuf = protocol.NewConnBuf()
	cBuf.Write(CH_alpn_h2)
	rule, err = tlsProtocol.Identify(ctx, cBuf)
	if err != nil {
		t.Errorf("Error identifying rule: %s", err)
	}
	if rule != "ALPN h2" {
		t.Errorf("Wrong rule identified: %s", rule)
	}

	cBuf = protocol.NewConnBuf()
	cBuf.Write(CH_catchall)
	rule, err = tlsProtocol.Identify(ctx, cBuf)
	if err != nil {
		t.Errorf("Error identifying rule: %s", err)
	}
	if rule != "CATCHALL" {
		t.Errorf("Wrong rule identified: %s", rule)
	}
}

func testIdentifyExpired(t *testing.T) {
	ctxExpired, cancel := context.WithCancel(context.Background())
	cancel()

	cBuf := protocol.NewConnBuf()
	cBuf.Write(CH_cloudflare_dns_com)
	rule, err := tlsProtocol.Identify(ctxExpired, cBuf)
	if err == nil {
		t.Errorf("should have returned error")
	}
	if rule != "" {
		t.Errorf("should have returned empty rule")
	}
}

func testIdentifyWithNotEnoughData(t *testing.T) {
	ctx := context.Background()
	ctxFiveSeconds, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cBuf := protocol.NewConnBuf()
	cBuf.Write(CH_cloudflare_dns_com[:5])
	rule, err := tlsProtocol.Identify(ctxFiveSeconds, cBuf)
	if err == nil {
		t.Errorf("should have returned error")
	}
	if rule != "" {
		t.Errorf("should have returned empty rule")
	}
}
