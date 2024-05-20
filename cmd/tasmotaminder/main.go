package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"tasmotamanager/types"
	"tasmotamanager/utils"
)

func main() {
	s := types.NewState(getRuleConfig(), notify)
	defer s.Stop()

	clientOptions := getClientOptions()

	clientOptions.SetDefaultPublishHandler(defaultReceiveHandler)
	clientOptions.OnConnect = getConnectedHandler(s)
	clientOptions.OnConnectionLost = getDisconnectedHandler(s)

	client := mqtt.NewClient(clientOptions)

	s.SetupSchedules(client)

	utils.WaitForToken(client.Connect())

	utils.Wait()
}
