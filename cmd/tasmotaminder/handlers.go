package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"tasmotamanager/types"
	"tasmotamanager/utils"
)

func defaultReceiveHandler(client mqtt.Client, msg mqtt.Message) {
	log.Println("Received message:", string(msg.Payload()), "from topic:", msg.Topic())
}

func connectHandler(client mqtt.Client) {
	log.Println("Connected")
}

func connectLostHandler(client mqtt.Client, err error) {
	log.Println("Disconnected", err)
}

func newConfigReceiveHandler(s *state) func(client mqtt.Client, msg mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		cfg := types.NewConfig()
		if err := json.Unmarshal(msg.Payload(), cfg); err != nil {
			log.Fatalf("error unmarshalling config: %s", err)
		}

		topic, topicErr := cfg.Topic(types.TeleSuffix, types.SensorSuffix)
		if topicErr != nil {
			log.Fatalf("error retrieving telemetry topic from config: %s", topicErr)
		}

		log.Println("Received config for Tasmota plug", cfg.Name())

		if !s.updateConfig(cfg) {
			utils.WaitForToken(client.Subscribe(topic, 2, getSensorHandler(s, cfg)))

			log.Println("Started new subscription for", topic)
		}
	}
}

func getSensorHandler(s *state, cfg *types.Config) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		sensor := types.NewSensor()
		if err := json.Unmarshal(msg.Payload(), sensor); err != nil {
			log.Fatalf("error unmarshalling sensor: %s", err)
		}

		log.Println("Received sensor message", cfg.Name())

		s.updateSensor(cfg, sensor)
	}
}
