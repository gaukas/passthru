package tls

import (
	"fmt"
	"strings"

	"github.com/gaukas/passthru/config"
)

func ValidateRule(rule string) error {
	// split into parts delimited by space
	ruleParts := strings.Split(rule, " ")
	if len(ruleParts) > 2 || len(ruleParts) < 1 {
		return fmt.Errorf("invalid rule: %s", rule)
	}

	// validate rule
	switch ruleParts[0] {
	case "SNI":
		if len(ruleParts) != 2 {
			return fmt.Errorf("invalid rule: %s", rule)
		}
	case "ALPN":
		if len(ruleParts) != 2 {
			return fmt.Errorf("invalid rule: %s", rule)
		}
	case "CATCHALL":
		if len(ruleParts) != 1 {
			return fmt.Errorf("invalid rule: %s", rule)
		}
	default:
		return fmt.Errorf("invalid rule: %s", rule)
	}

	return nil
}

// Rule type
const (
	RuleSNI uint8 = iota
	RuleALPN
	RuleCATCHALL
)

type Rule struct {
	Type     uint8
	Contents string
	RuleName config.Rule
}

func ParseRule(rule config.Rule) (Rule, error) {
	// validate rule
	err := ValidateRule(rule)
	if err != nil {
		return Rule{}, err
	}

	// split into parts delimited by space
	ruleParts := strings.Split(rule, " ")

	// parse rule
	switch ruleParts[0] {
	case "SNI":
		return Rule{
			Type:     RuleSNI,
			Contents: ruleParts[1],
			RuleName: rule,
		}, nil
	case "ALPN":
		return Rule{
			Type:     RuleALPN,
			Contents: ruleParts[1],
			RuleName: rule,
		}, nil
	case "CATCHALL":
		return Rule{
			Type:     RuleCATCHALL,
			RuleName: rule,
		}, nil
	default:
		return Rule{}, fmt.Errorf("invalid rule: %s", rule)
	}
}

func ParseRules(rules []config.Rule) ([]Rule, error) {
	var catchAllRule Rule

	parsedRules := []Rule{}
	for _, rule := range rules {
		parsedRule, err := ParseRule(rule)
		if err != nil {
			return []Rule{}, err
		}
		if parsedRule.Type == RuleCATCHALL {
			catchAllRule = parsedRule // catch all rule must be last
		} else {
			parsedRules = append(parsedRules, parsedRule)
		}
	}

	if catchAllRule.RuleName != "" {
		parsedRules = append(parsedRules, catchAllRule)
	}

	return parsedRules, nil
}
