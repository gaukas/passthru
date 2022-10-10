package protocol

import (
	"github.com/gaukas/passthru/config"
)

type Protocol interface {
	Name() config.Protocol                // Name of the protocol
	ApplyRules(rules []config.Rule) error // Apply rules to the protocol for later identification
	Identify(cBuf *ConnBuf) (config.Rule, error)
}
