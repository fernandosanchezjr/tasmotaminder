package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"tasmotamanager/utils"
)

func main() {
	s := newState()

	clientOptions := getClientOptions()

	clientOptions.SetDefaultPublishHandler(defaultReceiveHandler)
	clientOptions.OnConnect = getConnectedHandler(s)
	clientOptions.OnConnectionLost = disconnectedHandler

	client := mqtt.NewClient(clientOptions)
	utils.WaitForToken(client.Connect())

	utils.Wait()
}
