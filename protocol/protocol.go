package protocol

import (
	"context"

	"github.com/gaukas/passthru/config"
)

type Protocol interface {
	Name() config.Protocol                                            // Name of the protocol
	ApplyRules(rules []config.Rule) error                             // Apply rules to the protocol for later Identify()
	Identify(ctx context.Context, cBuf *ConnBuf) (config.Rule, error) // Identify will keep checking cBuf until it can identify a rule
}
