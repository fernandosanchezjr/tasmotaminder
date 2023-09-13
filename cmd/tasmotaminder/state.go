package main

import (
	"log"
	"tasmotamanager/types"
)

type plugState struct {
	sensor *types.Sensor
	state  *types.State
}

type state struct {
	plugs map[string]*plugState
}

func newPlugState() *plugState {
	return &plugState{}
}

func (ps *plugState) updateSensor(sensor *types.Sensor) {
	ps.sensor = sensor
}

func (ps *plugState) updateState(state *types.State) {
	ps.state = state
}

func newState() *state {
	return &state{plugs: make(map[string]*plugState)}
}

func (s *state) getOrCreatePlug(id string) *plugState {
	ps, found := s.plugs[id]
	if !found {
		ps = newPlugState()
		s.plugs[id] = ps
	}

	return ps
}

func (s *state) updateSensor(id string, sensor *types.Sensor) {
	ps := s.getOrCreatePlug(id)
	ps.updateSensor(sensor)

	log.Println("Updating sensor", id, "with", sensor)
}

func (s *state) updateState(id string, state *types.State) {
	ps := s.getOrCreatePlug(id)
	ps.updateState(state)

	log.Println("Updating state", id, "with", state)
}
