package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"tasmotamanager/types"
	"tasmotamanager/utils"
)

func getSubscriptionTopics() map[string]byte {
	return map[string]byte{
		types.TasmotaSensorTopic: 0,
		types.TasmotaStateTopic:  0,
	}
}

func defaultReceiveHandler(_ mqtt.Client, msg mqtt.Message) {
	log.Println("Received message:", string(msg.Payload()), "from topic:", msg.Topic())
}

func getConnectedHandler(s *types.State) mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		log.Println("Connected")

		utils.WaitForToken(client.SubscribeMultiple(
			getSubscriptionTopics(),
			getSensorHandler(s),
		))
	}
}

func disconnectedHandler(_ mqtt.Client, err error) {
	log.Println("Disconnected", err)
}

func getSensorHandler(s *types.State) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		deviceId := utils.GetDeviceId(topic)
		topicSuffix := utils.GetTopicSuffix(topic)
		payload := msg.Payload()
		var sensorMessage *types.Sensor
		var stateMessage *types.SensorState

		switch topicSuffix {
		case types.TasmotaSensorSuffix:
			sensorMessage = types.NewSensor()
			sensorMessage.Unmarshal(payload)
		case types.TasmotaStateSuffix:
			stateMessage = types.NewSensorState()
			stateMessage.Unmarshal(payload)
		}

		s.Update(client, deviceId, sensorMessage, stateMessage)
	}
}
