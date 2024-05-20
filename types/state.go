package types

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

type NotifyFunc func(title string, tags string, body string) error

type State struct {
	plugRules PlugRules
	plugs     map[string]*PlugState
	scheduler *gocron.Scheduler
	notify    NotifyFunc
}

func NewState(plugRules PlugRules, notifier NotifyFunc) *State {
	log.Printf("Loaded rules:\n%s", plugRules)

	s := &State{
		plugRules: plugRules,
		plugs:     make(map[string]*PlugState),
		scheduler: gocron.NewScheduler(time.Local),
		notify:    notifier, // Initialize the notify field
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
	rule *PlugRule,
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

		action.Execute(ps.getRuleTarget(client, rule))
		if rule.Notify {
			title := fmt.Sprintf("%s %s", rule.DeviceName(), action)
			body := fmt.Sprintf("Executing from schedule %s", cron)
			err := s.notify(title, "tasmotaminder,timer", body)
			if err != nil {
				log.Printf("Error sending notification: %v", err)
			}
		}
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
