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
	clientOptions.OnConnect = connectHandler
	clientOptions.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(clientOptions)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	utils.WaitForToken(client.Subscribe(types.TasmotaConfigTopic, 2, newConfigReceiveHandler(s)))

	utils.Wait()
}
