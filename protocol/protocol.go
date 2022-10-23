package protocol

import (
	"context"

	"github.com/gaukas/passthru/config"
)

type Protocol interface {
	// Name prints the name of the protocol, like "TLS", which is going to be used as a key in the ProtocolGroup
	Name() config.Protocol

	// ApplyRules inputs the rules to be used by the protocol.
	// Note the rules are out-of-order.
	// TODO: Input in-order rules instead of out-of-order rules. May need something other than a map.
	// Protocol implementations should make sure the CATCHEALL rule is always the last rule to be applied.
	ApplyRules(rules []config.Rule) error // Apply rules to the protocol for later Identify()

	// Identify identifies the rule that matches the request.
	Identify(ctx context.Context, cBuf *ConnBuf) (config.Rule, error) // Identify will keep checking cBuf until it can identify a rule
}
