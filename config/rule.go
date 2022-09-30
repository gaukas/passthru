package config

// Example Action:
// {
// 		"action": "FORWARD",
// 		"to_addr": "gaukas.wang:443"
// }

// Action is a struct representing an action to be taken
// on a request that matches a rule
type Action struct {
	Type   ActionType `json:"type"`    // Type of action to take when the rule is matched
	ToAddr string     `json:"to_addr"` // Address to FORWARD to, if type is FORWARD
}

type ActionType = uint8

const (
	ACTION_FORWARD ActionType = iota + 1 // "FORWARD"
	ACTION_REJECT                        // "REJECT"
)

// Implement custom unmarshaller/marshaller for ActionType
// Due to type conflict. (JSON: string, Go: uint8)

func (at *ActionType) UnmarshalJSON(data []byte) error {
	return nil
}

func (at *ActionType) MarshalJSON() ([]byte, error) {
	return nil, nil
}
