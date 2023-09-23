package types

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

const (
	Off RuleAction = iota
	On
	Reset
)

type RuleAction int

func (ra *RuleAction) String() string {
	if ra == nil {
		return ""
	}

	switch *ra {
	case Off:
		return "off"
	case On:
		return "on"
	case Reset:
		return "reset"
	}

	return ""
}

func (ra *RuleAction) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: ra.String(),
	}, nil
}

func (ra *RuleAction) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return fmt.Errorf("invalid node: %+#v", node)
	}
	switch strings.ToLower(node.Value) {
	case "off":
		*ra = Off
	case "on":
		*ra = On
	case "reset":
		*ra = Reset
	default:
		return fmt.Errorf("invalid node value: %s", node.Value)
	}

	return nil
}

func (ra *RuleAction) Execute(target RuleTarget) {
	if ra == nil || *ra == Off {
		target.Off()
		return
	}

	switch *ra {
	case On:
		target.On()
	case Reset:
		target.Reset()
	}
}
