package protocol

import (
	"io"

	"github.com/gaukas/passthru/config"
)

type ProtocolManager struct {
	protocols map[config.Protocol]Protocol
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

}

func (pm *ProtocolManager) FindAction(conn io.Reader) (config.Action, error) {

}
