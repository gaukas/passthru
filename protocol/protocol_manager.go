package protocol

import (
	"context"
	"fmt"

	"github.com/gaukas/passthru/config"
)

type ProtocolManager struct {
	protocols     map[config.Protocol]Protocol
	protocolGroup config.ProtocolGroup
}

func NewProtocolManager() *ProtocolManager {
	return &ProtocolManager{
		protocols: make(map[config.Protocol]Protocol),
	}
}

// Called before ImportProtocolGroup, or will see error upon unknown protocol
func (pm *ProtocolManager) RegisterProtocol(p Protocol) {
	pm.protocols[p.Name()] = p
}

func (pm *ProtocolManager) GetProtocol(name config.Protocol) Protocol {
	return pm.protocols[name]
}

// Called after RegisterProtocol, or will see error upon unknown protocol
func (pm *ProtocolManager) ImportProtocolGroup(pg config.ProtocolGroup) error {
	for protocol, filter := range pg {
		p := pm.GetProtocol(protocol)
		if p == nil {
			return fmt.Errorf("unknown protocol: %s", protocol)
		}
		rules := []config.Rule{}
		for rule := range filter {
			rules = append(rules, rule)
		}
		err := p.ApplyRules(rules)
		if err != nil {
			return err
		}
	}

	pm.protocolGroup = pg

	return nil
}

func (pm *ProtocolManager) FindAction(ctx context.Context, cBuf *ConnBuf) (config.Action, error) {
	chanRule := make(chan config.Rule)
	chanProtocol := make(chan config.Protocol)
	subctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for pName, p := range pm.protocols {
		go func(protocolName config.Protocol, protocol Protocol) {
			rule, err := protocol.Identify(subctx, cBuf)
			if err != nil {
				return
			}
			chanRule <- rule
			chanProtocol <- protocolName
		}(pName, p)
	}

	select {
	case rule := <-chanRule:
		select {
		case protocolName := <-chanProtocol:
			// look for the rule in the protocol group
			filter, ok := pm.protocolGroup[protocolName]
			if !ok {
				return config.Action{}, fmt.Errorf("unknown protocol: %s", protocolName)
			}

			action, ok := filter[rule]
			if !ok {
				return config.Action{}, fmt.Errorf("unknown rule: %s", rule)
			}
			return action, nil
		case <-ctx.Done():
			return config.Action{}, ctx.Err()
		}
	case <-ctx.Done():
		return config.Action{}, ctx.Err()
	}
}
