package types

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

const (
	EqualTo PowerComparison = iota
	LessThan
	GreaterThan
)

type PowerComparison int

func (pc *PowerComparison) String() string {
	if pc == nil {
		return "greaterThan"
	}

	switch *pc {
	case EqualTo:
		return "equalTo"
	case LessThan:
		return "lessThan"
	case GreaterThan:
		return "greaterThan"
	}

	return ""
}

func (pc *PowerComparison) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: pc.String(),
	}, nil
}

func (pc *PowerComparison) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return fmt.Errorf("invalid node: %+#v", node)
	}
	switch strings.ToLower(node.Value) {
	case "equalto":
		*pc = EqualTo
	case "lessthan":
		*pc = LessThan
	case "greaterthan":
		*pc = GreaterThan
	default:
		return fmt.Errorf("invalid node value: %s", node.Value)
	}

	return nil
}

func (pc *PowerComparison) Compare(plug int, target int) bool {
	if pc == nil || *pc == GreaterThan {
		return plug > target
	}
	switch *pc {
	case EqualTo:
		return plug == target
	case LessThan:
		return plug < target
	}

	return false
}
