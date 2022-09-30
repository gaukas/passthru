package config

type Server struct {
	ListenOn string `json:"listen_on"` // An address to listen on. E.g.: 0.0.0.0:443

	// A list of filters to apply to the request
	// key: protocol name
	// value: filter (including rules)
	Filters map[string]Filter `json:"filters"`
}

// Shouldn't need to implement custom unmarshaller/marshaller for Server
