package protocol_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gaukas/passthru/config"
	"github.com/gaukas/passthru/protocol"
)

type DummyProtocol struct {
}

func (p *DummyProtocol) Name() config.Protocol {
	return config.Protocol("dummy")
}

func (p *DummyProtocol) ApplyRules(rules []config.Rule) error {
	for _, rule := range rules {
		fmt.Printf("Applying rule: %s\n", rule)
	}
	return nil
}

func (p *DummyProtocol) Identify(ctx context.Context, cBuf *protocol.ConnBuf) (config.Rule, error) {
	if ctx.Err() != nil {
		return config.Rule(""), ctx.Err()
	}

	buf := make([]byte, 4)

	err := cBuf.Peek(buf, 4)
	if err != nil {
		return config.Rule(""), err
	}

	if string(buf) == "test" {
		return config.Rule("test"), nil
	}

	return config.Rule(""), errors.New("no match rules")
}

var (
	pm *protocol.ProtocolManager
)

func TestProtocolManager(t *testing.T) {
	testNewProtocolManager(t)
	testRegisterProtocol(t)
	testGetProtocol(t)
	testImportProtocolGroup(t)
	testFindAction(t)
}

func testNewProtocolManager(t *testing.T) {
	pm = protocol.NewProtocolManager()
	if pm == nil {
		t.Errorf("Error creating ProtocolManager")
	}
}

func testRegisterProtocol(t *testing.T) {
	pm.RegisterProtocol(&DummyProtocol{})
}

func testGetProtocol(t *testing.T) {
	p := pm.GetProtocol(config.Protocol("dummy"))
	if p == nil {
		t.Errorf("Error getting protocol")
	}
}

func testImportProtocolGroup(t *testing.T) {
	pg := config.ProtocolGroup{
		config.Protocol("dummy"): config.Filter{
			config.Rule("test"): config.Action{
				Type:   config.ACTION_FORWARD,
				ToAddr: "127.0.0.2:8080",
			},
		},
	}

	err := pm.ImportProtocolGroup(pg)
	if err != nil {
		t.Errorf("Error importing protocol group: %s", err)
	}
}

func testFindAction(t *testing.T) {
	cBuf := protocol.NewConnBuf()
	cBuf.Write([]byte("test"))
	action, err := pm.FindAction(context.Background(), cBuf)
	if err != nil {
		t.Errorf("Error finding action: %s", err)
	}

	if action.Type != config.ACTION_FORWARD {
		t.Errorf("Error finding action: %s", err)
	}
}
