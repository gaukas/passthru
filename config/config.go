package config

import (
	"encoding/json"
	"os"

	"github.com/gaukas/passthru/internal/logger"
)

// Config is a struct that can be loaded from a JSON file
// or written to a JSON file
type Config struct {
	Version Version     `json:"version"`
	Servers ServerGroup `json:"servers"` // A list of servers to listen on
}

func LoadConfig(filename string) (*Config, error) {
	// read data from file
	logger.Debugf("Loading config from %s", filename)
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// then call json unmarshal
	c := Config{}
	json.Unmarshal(content, &c)

	return &c, nil
}

func (c *Config) Write(filename string) error {
	// call json marshal
	content, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// then write data to file
	logger.Debugf("Writing config to %s", filename)
	err = os.WriteFile(filename, content, 0644)
	return err
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
