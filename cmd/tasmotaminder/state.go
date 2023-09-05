package main

import (
	"log"
	"tasmotamanager/types"
)

type plugState struct {
	config *types.Config
	sensor *types.Sensor
}

type state struct {
	plugs map[string]*plugState
}

func newState() *state {
	return &state{plugs: make(map[string]*plugState)}
}

func newPlugState(config *types.Config) *plugState {
	return &plugState{config: config}
}

func (s *state) updateConfig(cfg *types.Config) bool {
	plug, found := s.plugs[cfg.MAC]
	if !found {
		s.plugs[cfg.MAC] = newPlugState(cfg)
	} else {
		plug.config = cfg
	}

	return found
}

func (s *state) updateSensor(cfg *types.Config, sensor *types.Sensor) {
	plug, found := s.plugs[cfg.MAC]
	if found {
		plug.sensor = sensor
		s.stateUpdate(plug)
	} else {
		log.Println("during sensor update, could not find plug", cfg.MAC)
	}
}

func (s *state) stateUpdate(plug *plugState) {
	log.Println("State update event for", plug.config)
}
