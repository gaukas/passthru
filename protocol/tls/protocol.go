package tls

import (
	"context"
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

func (p *Protocol) Clone() protocol.Protocol {
	pCopy := &Protocol{}
	for _, rule := range p.rules {
		pCopy.rules = append(pCopy.rules, rule)
	}
	return pCopy
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

func (p *Protocol) Identify(ctx context.Context, cBuf *protocol.ConnBuf) (config.Rule, error) {
	connInfo, err := ParseClientHello(ctx, cBuf)
	if err != nil {
		return "", err
	}

	// identify rule by the original order
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
