package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"tasmotamanager/types"
	"tasmotamanager/utils"
)

func main() {
	s := types.NewState(getRuleConfig())

	clientOptions := getClientOptions()

	clientOptions.SetDefaultPublishHandler(defaultReceiveHandler)
	clientOptions.OnConnect = getConnectedHandler(s)
	clientOptions.OnConnectionLost = disconnectedHandler

	client := mqtt.NewClient(clientOptions)
	utils.WaitForToken(client.Connect())

	utils.Wait()
}
