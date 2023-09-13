package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"tasmotamanager/types"
	"tasmotamanager/utils"
)

func main() {
	s := newState()

	clientOptions := getClientOptions()

	clientOptions.SetDefaultPublishHandler(defaultReceiveHandler)
	clientOptions.OnConnect = connectedHandler
	clientOptions.OnConnectionLost = disconnectedHandler

	client := mqtt.NewClient(clientOptions)
	utils.WaitForToken(client.Connect())

	utils.WaitForToken(client.SubscribeMultiple(
		map[string]byte{
			types.TasmotaSensorTopic: 1,
			types.TasmotaStateTopic:  1,
		},
		getSensorHandler(s),
	))

	utils.Wait()
}
