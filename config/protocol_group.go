package config

// Example ProtocolGroup:
// {
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
// }

// ProtocolGroup includes filters to apply to the request per each protocol
type ProtocolGroup = map[Protocol]Filter

// Protocol is a string representing the protocol name
type Protocol = string // E.g.: "TLS", "SSH"
