package main

import (
	"log"
	"tasmotamanager/types"
	"time"
)

const maxMessageTimeDifference = time.Second * 10

type plugState struct {
	id            string
	sensor        *types.Sensor
	state         *types.State
	sensorUpdated time.Time
	stateUpdated  time.Time
}

type state struct {
	plugs map[string]*plugState
}

func newPlugState(id string) *plugState {
	return &plugState{id: id}
}

func (ps *plugState) resetTime() {
	ps.sensorUpdated = time.Time{}
	ps.stateUpdated = time.Time{}
}

func (ps *plugState) triggerEvent() {
	if ps.sensorUpdated.IsZero() || ps.stateUpdated.IsZero() {
		return
	}

	difference := ps.sensorUpdated.Sub(ps.stateUpdated).Abs()
	ps.resetTime()

	if difference > maxMessageTimeDifference {
		log.Println("update time difference too large:", difference)
	}

	log.Println(ps.id, "update message latency", difference)
}

func (ps *plugState) updateSensor(sensor *types.Sensor) {
	ps.sensor = sensor
	ps.sensorUpdated = time.Now()

	log.Println("Updating sensor", ps.id)

	ps.triggerEvent()
}

func (ps *plugState) updateState(state *types.State) {
	ps.state = state
	ps.stateUpdated = time.Now()

	log.Println("Updating state", ps.id)

	ps.triggerEvent()
}

func newState() *state {
	return &state{plugs: make(map[string]*plugState)}
}

func (s *state) getOrCreatePlug(id string) *plugState {
	ps, found := s.plugs[id]
	if !found {
		ps = newPlugState(id)
		s.plugs[id] = ps
	}

	return ps
}

func (s *state) update(id string, sensor *types.Sensor, state *types.State) {
	ps := s.getOrCreatePlug(id)

	if sensor != nil {
		ps.updateSensor(sensor)
	}

	if state != nil {
		ps.updateState(state)
	}
}
