package types

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Config struct {
	IP    string         `json:"ip"`
	DN    string         `json:"dn"`
	FN    []interface{}  `json:"fn"`
	HN    string         `json:"hn"`
	MAC   string         `json:"mac"`
	MD    string         `json:"md"`
	TY    int            `json:"ty"`
	IF    int            `json:"if"`
	OFLN  string         `json:"ofln"`
	ONLN  string         `json:"onln"`
	STATE []string       `json:"state"`
	SW    string         `json:"sw"`
	T     string         `json:"t"`
	FT    string         `json:"ft"`
	TP    []string       `json:"tp"`
	RL    []int          `json:"rl"`
	SWC   []int          `json:"swc"`
	SWN   []interface{}  `json:"swn"`
	BTN   []int          `json:"btn"`
	SO    map[string]int `json:"so"`
	LK    int            `json:"lk"`
	LKST  int            `json:"lk_st"`
	SHO   []int          `json:"sho"`
	VER   int            `json:"ver"`
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) String() string {
	data, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("error marshalling Config to JSON: %s", err)
	}
	return string(data)
}

func (cfg *Config) HasTopicPrefix(topic string) bool {
	for _, tp := range cfg.TP {
		if topic == tp {
			return true
		}
	}
	return false
}

func (cfg *Config) Topic(prefix string, suffix string) (string, error) {
	if !cfg.HasTopicPrefix(prefix) {
		return "", fmt.Errorf("config has no prefix: %s", prefix)
	}

	components := []string{prefix, cfg.T}
	if suffix != "" {
		components = append(components, suffix)
	}

	return strings.Join(components, "/"), nil
}

func (cfg *Config) Name() string {
	return fmt.Sprintf("%s (%s, mac address %s)", cfg.DN, cfg.MD, cfg.MAC)
}
