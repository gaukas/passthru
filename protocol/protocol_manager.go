package protocol

import (
	"fmt"

	"github.com/gaukas/passthru/config"
)

type ProtocolManager struct {
	protocols     map[config.Protocol]Protocol
	protocolGroup config.ProtocolGroup
}

func NewProtocolManager() *ProtocolManager {
	return &ProtocolManager{}
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
			return fmt.Errorf("Unknown protocol: %s", protocol)
		}
		rules := []config.Rule{}
		for rule, _ := range filter {
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

func (pm *ProtocolManager) FindAction(cBuf *ConnBuf) (config.Action, error) {
	for pName, p := range pm.protocols {
		rule, err := p.Identify(cBuf)
		if err != nil {
			continue
		}

		// look for the rule in the protocol group
		filter, ok := pm.protocolGroup[pName]
		if !ok {
			return config.Action{}, fmt.Errorf("Unknown protocol: %s", pName)
		}

		action, ok := filter[rule]
		if !ok {
			return config.Action{}, fmt.Errorf("Unknown rule: %s", rule)
		}

		return action, nil
	}
	return config.Action{}, fmt.Errorf("No rule matched")
}
