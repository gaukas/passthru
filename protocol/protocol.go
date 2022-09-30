package protocol

import (
	"io"

	"github.com/gaukas/passthru/config"
)

type Protocol interface {
	Name() config.Protocol // Name of the protocol
	ApplyRules(rules []config.Rule) error
	Identify(conn io.Reader) (config.Rule, error)
}
