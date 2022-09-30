package config

// Example Filter:
// {
// 	"SNI gaukas.wang": {
// 		"action": "FORWARD",
// 		"to_addr": "gaukas.wang:443"
// 	},
// 	"SNI google.com": {
// 		"action": "FORWARD",
// 		"to_addr": "google.com:443"
// 	},
// 	"CATCHALL": {
// 		"action": "REJECT"
// 	}
// }

// Filter defines the Rule to Action mapping relationship.
type Filter = map[Rule]Action

// Rule is a string that can be matched against a request by a filter
type Rule = string // E.g.: "SNI example.com", "CATCHALL"
