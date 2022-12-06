package protocol

import (
	"context"

	"github.com/gaukas/passthru/config"
)

// Protocol is the interface for protocol identification.
type Protocol interface {
	// Name prints the name of the protocol, like "TLS", which is going to be used as a key in the ProtocolGroup
	Name() config.Protocol

	// Clone creates a new Protocol instance with the same rules (as a deep copy)
	Clone() Protocol

	// ApplyRules save the rules for later Identify calls.
	// Note the rules are out-of-order intentionally to prevent conflicting rules.
	// Protocol implementations should make sure the CATCHEALL rule is always the last rule to be applied.
	ApplyRules(rules []config.Rule) error

	// Identify identifies the rule that matches the request.
	Identify(ctx context.Context, cBuf *ConnBuf) (config.Rule, error) // Identify will keep checking cBuf until it can make a deterministic decision or the context is cancelled.
}
