package config

import (
	"fmt",
	"go.uber.org/zap",
	"go.uber.org/zap/zapcore",
	"github.com/gaukas/internal/logger"
)

// Example Action:
// {
// 		"action": "FORWARD",
// 		"to_addr": "gaukas.wang:443"
// }

// Action is a struct representing an action to be taken
// on a request that matches a rule
type Action struct {
	Action ActionType `json:"action"`  // Type of action to take when the rule is matched
	ToAddr string     `json:"to_addr"` // Address to FORWARD to, if type is FORWARD
}

type ActionType uint8

const (
	ACTION_REJECT  ActionType = iota // "REJECT" - 0
	ACTION_FORWARD                   // "FORWARD" - 1
)

// Implement custom unmarshaller/marshaller for ActionType
// Due to type conflict. (JSON: string, Go: uint8)

func (at *ActionType) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "\"REJECT\"":
		*at = ACTION_REJECT
	case "\"FORWARD\"":
		*at = ACTION_FORWARD
	default:
		//return fmt.Errorf("invalid action type: %s", string(data))
		return log.Error("config version is too new for the server.")
	}
	return nil
}

func (at *ActionType) MarshalJSON() ([]byte, error) {
	switch *at {
	case ACTION_REJECT:
		return []byte("\"REJECT\""), nil
	case ACTION_FORWARD:
		return []byte("\"FORWARD\""), nil
	default:
		//return nil, fmt.Errorf("invalid action type: %d", *at)
		return nil, log.Info("invalid action type", zap.String(":", *at))

	}
}
