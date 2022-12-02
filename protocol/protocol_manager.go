package protocol

import (
	"context"
	"fmt"
	"sync"

	"github.com/gaukas/passthru/config"
	"github.com/gaukas/passthru/internal/logger"
)

type ProtocolManager struct {
	protocols     map[config.Protocol]Protocol
	protocolGroup config.ProtocolGroup
	catchAll      config.Action
}

func NewProtocolManager() *ProtocolManager {
	return &ProtocolManager{
		protocols: make(map[config.Protocol]Protocol),
	}
}

// Called before ImportProtocolGroup, or will see error upon unknown protocol
func (pm *ProtocolManager) RegisterProtocol(p Protocol) {
	pm.protocols[p.Name()] = p.Clone()
}

func (pm *ProtocolManager) GetProtocol(name config.Protocol) Protocol {
	return pm.protocols[name]
}

// Called after RegisterProtocol, or will see error upon unknown protocol
func (pm *ProtocolManager) ImportProtocolGroup(pg config.ProtocolGroup) error {
	logger.Infof("Importing protocol group...")
LOOP_PG:
	for protocol, filter := range pg {
		logger.Debugf("Importing protocol %s", protocol)
		if protocol == "CATCHALL" { // if CATCHALL, save it in catchAll.
			for rule, action := range filter {
				if rule != "CATCHALL" {
					return fmt.Errorf("the CATCHALL protocol must ONLY have CATCHALL rule")
				}
				pm.catchAll = action
				continue LOOP_PG
			}
			return fmt.Errorf("the CATCHALL protocol must have CATCHALL rule")
		} // When not set, CATCHALL will be REJECT

		p := pm.GetProtocol(protocol)
		if p == nil {
			return fmt.Errorf("unknown protocol: %s", protocol)
		}
		rules := []config.Rule{}
		for rule := range filter {
			logger.Debugf("Importing rule %s", rule)
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

	wg := &sync.WaitGroup{}

	for pName, p := range pm.protocols {
		wg.Add(1)
		go func(protocolName config.Protocol, protocol Protocol) {
			defer wg.Done()
			rule, err := protocol.Identify(subctx, cBuf)
			if err != nil {
				return
			}
			chanRule <- rule
			chanProtocol <- protocolName
		}(pName, p)
	}

	// wait for all goroutines to finish, then cancel the context
	go func() {
		wg.Wait()
		cancel()
	}()

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

			logger.Debugf("Found action %s for protocol %s and rule %s", action, protocolName, rule)
			return action, nil
		case <-ctx.Done(): // CATCHALL
			return pm.catchAll, ctx.Err()
		}
	case <-ctx.Done():
		return pm.catchAll, ctx.Err()
	}
}
