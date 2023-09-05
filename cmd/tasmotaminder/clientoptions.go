package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"tasmotamanager/utils"
)

const (
	defaultHost     = "localhost"
	defaultPort     = "1883"
	defaultClientId = "tasmotaminder"
)

func getBrokerUrl() string {
	host := utils.GetEnvOrDefault("BROKER_HOST", defaultHost)
	port := utils.GetEnvOrDefault("BROKER_PORT", defaultPort)

	return fmt.Sprintf("tcp://%s:%s", host, port)
}

func getClientOptions() *mqtt.ClientOptions {
	clientOptions := mqtt.NewClientOptions()
	clientOptions.AddBroker(getBrokerUrl())
	clientOptions.SetClientID(utils.GetEnvOrDefault("CLIENT_ID", defaultClientId))
	clientOptions.SetUsername(os.Getenv("BROKER_USERNAME"))
	clientOptions.SetPassword(os.Getenv("BROKER_PASSWORD"))

	return clientOptions
}
