package config

type Rule struct {
	Action Action `json:"action"`  // Action to take when the rule is matched
	ToAddr string `json:"to_addr"` // Address to redirect to
}

type Action = uint8

const (
	ACTION_FORWARD Action = iota + 1
	ACTION_REJECT
)

// Implement custom unmarshaller/marshaller for Action
// Due to type conflict. (JSON: string, Go: uint8)

func (a *Action) UnmarshalJSON(data []byte) error {
	return nil
}

func (a *Action) MarshalJSON() ([]byte, error) {
	return nil, nil
}
