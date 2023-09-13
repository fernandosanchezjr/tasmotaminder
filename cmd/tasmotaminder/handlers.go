package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
	"tasmotamanager/types"
	"tasmotamanager/utils"
)

func defaultReceiveHandler(_ mqtt.Client, msg mqtt.Message) {
	log.Println("Received message:", string(msg.Payload()), "from topic:", msg.Topic())
}

func connectedHandler(_ mqtt.Client) {
	log.Println("Connected")
}

func disconnectedHandler(_ mqtt.Client, err error) {
	log.Println("Disconnected", err)
}

func getSensorHandler(s *state) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		deviceId := utils.GetDeviceId(topic)
		payload := msg.Payload()

		if strings.HasSuffix(topic, types.TasmotaSensorSuffix) {
			ps := types.NewSensor()
			ps.Unmarshal(payload)

			s.updateSensor(deviceId, ps)
		}

		if strings.HasSuffix(topic, types.TasmotaStateSuffix) {
			ps := types.NewState()
			ps.Unmarshal(payload)

			s.updateState(deviceId, ps)
		}
	}
}
