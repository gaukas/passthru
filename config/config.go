package config

// Config is a struct that can be loaded from a JSON file
// or written to a JSON file
type Config struct {
	MinVersion Version     `json:"min_version"` // TODO: implement marshaller/unmarshaller for Version
	Servers    ServerGroup `json:"servers"`     // A list of servers to listen on
}

func (c *Config) Load(filename string) error {
	// read data from file
	// then call json unmarshal
}

func (c *Config) Write(filename string) error {
	// call json marshal
	// then write data to file
}

// Example Servers:
// "0.0.0.0:443": {
// 	"TLS": {
// 		"SNI gaukas.wang": {
// 			"action": "FORWARD",
// 			"to_addr": "gaukas.wang:443"
// 		},
// 		"SNI google.com": {
// 			"action": "FORWARD",
// 			"to_addr": "google.com:443"
// 		},
// 		"CATCHALL": {
// 			"action": "REJECT"
// 		}
// 	},
// 	"SSH": {
// 		"CATCHALL": {
// 			"action": "REJECT"
// 		}
// 	}
// },
// "0.0.0.0:22": {
// 	"SSH": {
// 		"rules": {
// 			"CATCHALL": {
// 				"action": "FORWARD",
// 				"to_addr": "127.0.0.1:22122"
// 			}
// 		}
// 	}
// }

// ServerGroup is a map of server address to protocol filters
type ServerGroup = map[ServerAddr]ProtocolGroup
type ServerAddr = string
