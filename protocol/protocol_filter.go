package protocol

import "github.com/gaukas/passthru/config"

type ProtocolFilter struct {
	Protocol Protocol
	Filter   config.Filter
}
