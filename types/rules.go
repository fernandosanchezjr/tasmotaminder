package types

import (
	"gopkg.in/yaml.v3"
	"log"
)

type PlugRule struct {
	DeviceId             string           `yaml:"deviceId"`
	ResetDurationSeconds int              `yaml:"resetDurationSeconds,omitempty"`
	PowerTimer           *PowerTimer      `yaml:"powerTimer,omitempty"`
	PowerSchedules       []*PowerSchedule `yaml:"powerSchedules,omitempty"`
	Notify               bool             `yaml:"notify,omitempty"`   // Added Notify field
	Nickname             string           `yaml:"nickname,omitempty"` // Added Nickname field
}

type PlugRules []*PlugRule

func (pr PlugRules) String() string {
	data, err := yaml.Marshal(pr)
	if err != nil {
		log.Fatalf("error marshalling PlugRules to YAML: %s", err)
	}

	return string(data)
}

func (pr PlugRules) GetPlug(deviceId string) *PlugRule {
	for _, p := range pr {
		if p.DeviceId == deviceId {
			return p
		}
	}

	return nil
}

func (p *PlugRule) String() string {
	data, err := yaml.Marshal(p)
	if err != nil {
		log.Fatalf("error marshalling PlugRule to YAML: %s", err)
	}

	return string(data)
}

func (p *PlugRule) Evaluate(state *State, plug *PlugState, target RuleTarget) {
	if p.PowerTimer != nil {
		p.PowerTimer.Evaluate(state, plug, p, target)
	}
}

// DeviceName returns the Nickname if it's not an empty string, otherwise returns the DeviceId
func (p *PlugRule) DeviceName() string {
	if p.Nickname != "" {
		return p.Nickname
	}
	return p.DeviceId
}
