package tls

import (
	"errors"

	"github.com/gaukas/passthru/config"
	"github.com/gaukas/passthru/protocol"
)

type Protocol struct {
	rules []Rule
}

func (p *Protocol) Name() config.Protocol {
	return "TLS"
}

func (p *Protocol) ApplyRules(rules []config.Rule) error {
	// parse rules
	parsedRules, err := ParseRules(rules)
	if err != nil {
		return err
	}

	p.rules = parsedRules

	return nil
}

func (p *Protocol) Identify(cBuf *protocol.ConnBuf) (config.Rule, error) {
	connInfo, err := ParseConnInfo(cBuf)
	if err != nil {
		return "", err
	}

	// identify rule
	for _, rule := range p.rules {
		switch rule.Type {
		case RuleSNI:
			if connInfo.SNI == rule.Contents {
				return rule.RuleName, nil
			}
		case RuleALPN:
			if connInfo.ALPN == rule.Contents {
				return rule.RuleName, nil
			}
		case RuleCATCHALL:
			return rule.RuleName, nil
		}
	}

	return "", errors.New("no rule matched")
}
