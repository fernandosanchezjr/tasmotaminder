package types

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

type State struct {
	plugRules PlugRules
	plugs     map[string]*PlugState
	scheduler *gocron.Scheduler
}

func NewState(plugRules PlugRules) *State {
	log.Printf("Loaded rules:\n%s", plugRules)

	s := &State{
		plugRules: plugRules,
		plugs:     make(map[string]*PlugState),
		scheduler: gocron.NewScheduler(time.Local),
	}

	return s
}

func (s *State) SetupSchedules(client mqtt.Client) {
	for _, plugRule := range s.plugRules {
		deviceId := plugRule.DeviceId

		for _, schedule := range plugRule.PowerSchedules {

			cron, action := schedule.Cron, schedule.Action

			_, err := s.scheduler.Cron(cron).Do(
				s.getScheduleHandler(client, deviceId, cron, action, plugRule))

			if err != nil {
				log.Fatal("Invalid cron", err)
			}
			log.Println(
				"Started schedule",
				schedule.Cron,
				schedule.Action,
				"for plug",
				plugRule.DeviceId,
			)
		}
	}
}

func (s *State) getScheduleHandler(
	client mqtt.Client,
	deviceId string,
	cron string,
	action *RuleAction,
	plugRule *PlugRule,
) func() {
	return func() {
		ps := s.getOrCreatePlug(deviceId)

		log.Println(
			"Executing on schedule",
			cron,
			action,
			"for plug",
			deviceId,
		)

		action.Execute(ps.getRuleTarget(client, plugRule))
	}
}

func (s *State) getOrCreatePlug(id string) *PlugState {
	ps, found := s.plugs[id]
	if !found {
		ps = newPlugState(id)
		s.plugs[id] = ps
	}

	return ps
}

func (s *State) Update(client mqtt.Client, id string, sensor *Sensor, state *SensorState) {
	ps := s.getOrCreatePlug(id)

	if sensor != nil {
		ps.updateSensor(client, s, sensor)
	}

	if state != nil {
		ps.updateSensorState(client, s, state)
	}
}

func (s *State) Start() {
	log.Println("Starting scheduler")
	s.scheduler.StartAsync()
}

func (s *State) Stop() {
	log.Println("Stopping scheduler")
	s.scheduler.Stop()
}
