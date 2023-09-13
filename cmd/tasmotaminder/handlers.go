package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"tasmotamanager/types"
	"tasmotamanager/utils"
)

func defaultReceiveHandler(_ mqtt.Client, msg mqtt.Message) {
	log.Println("Received message:", string(msg.Payload()), "from topic:", msg.Topic())
}

func getConnectedHandler(s *state) mqtt.OnConnectHandler {
	log.Println("Connected")

	return func(client mqtt.Client) {
		utils.WaitForToken(client.SubscribeMultiple(
			map[string]byte{
				types.TasmotaSensorTopic: 1,
				types.TasmotaStateTopic:  1,
			},
			getSensorHandler(s),
		))
	}
}

func disconnectedHandler(_ mqtt.Client, err error) {
	log.Println("Disconnected", err)
}

func getSensorHandler(s *state) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		deviceId := utils.GetDeviceId(topic)
		topicSuffix := utils.GetTopicSuffix(topic)
		payload := msg.Payload()
		var sensorMessage *types.Sensor
		var stateMessage *types.State

		switch topicSuffix {
		case types.TasmotaSensorSuffix:
			sensorMessage = types.NewSensor()
			sensorMessage.Unmarshal(payload)
		case types.TasmotaStateSuffix:
			stateMessage = types.NewState()
			stateMessage.Unmarshal(payload)
		}

		s.update(deviceId, sensorMessage, stateMessage)

	}
}
