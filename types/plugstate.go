package types

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
	"tasmotamanager/utils"
	"time"
)

const maxMessageTimeDifference = time.Second * 10

type PlugState struct {
	id                 string
	sensor             *Sensor
	sensorState        *SensorState
	sensorUpdated      time.Time
	sensorStateUpdated time.Time
}

type plugClient struct {
	resetDuration time.Duration
	plug          *PlugState
	client        mqtt.Client
}

func newPlugState(id string) *PlugState {
	return &PlugState{id: id}
}

func (ps *PlugState) resetTime() {
	ps.sensorUpdated = time.Time{}
	ps.sensorStateUpdated = time.Time{}
}

func (ps *PlugState) getRuleTarget(client mqtt.Client, rule *PlugRule) RuleTarget {
	return &plugClient{
		resetDuration: utils.DurationMax(
			time.Second*time.Duration(rule.ResetDurationSeconds),
			TasmotaDefaultResetDuration,
		),
		plug:   ps,
		client: client,
	}
}

func (ps *PlugState) triggerEvent(client mqtt.Client, s *State) {
	if ps.sensorUpdated.IsZero() || ps.sensorStateUpdated.IsZero() {
		return
	}

	difference := ps.sensorUpdated.Sub(ps.sensorStateUpdated).Abs()
	ps.resetTime()

	if difference > maxMessageTimeDifference {
		log.Println("update time difference too large:", difference)
	}

	log.Println("Updated", ps.id)

	rule := s.plugRules.GetPlug(ps.id)
	if rule == nil {
		return
	}

	log.Printf("Evaluating rule:\n%s", rule)
	go rule.Evaluate(s, ps, ps.getRuleTarget(client, rule))
}

func (ps *PlugState) updateSensor(client mqtt.Client, state *State, sensor *Sensor) {
	ps.sensor = sensor
	ps.sensorUpdated = time.Now()

	ps.triggerEvent(client, state)
}

func (ps *PlugState) updateSensorState(client mqtt.Client, state *State, sensorState *SensorState) {
	ps.sensorState = sensorState
	ps.sensorStateUpdated = time.Now()

	ps.triggerEvent(client, state)
}

func (pc *plugClient) topic() string {
	return strings.ReplaceAll(TasmotaCommandTopic, TasmotaCommandWildcard, pc.plug.id)
}

func (pc *plugClient) publish(value string) {
	topic := pc.topic()
	log.Println("Publishing to", topic, value)
	utils.WaitForToken(pc.client.Publish(topic, 2, true, value))
}

func (pc *plugClient) Off() {
	log.Println("Turning off", pc.plug.id)
	pc.publish(TasmotaPowerOFF)
}

func (pc *plugClient) On() {
	log.Println("Turning on", pc.plug.id)
	pc.publish(TasmotaPowerON)
}

func (pc *plugClient) Reset() {
	log.Println("Resetting", pc.plug.id)
	pc.Off()

	time.Sleep(pc.resetDuration)

	pc.On()
}
