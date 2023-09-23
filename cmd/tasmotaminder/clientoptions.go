package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"tasmotamanager/types"
	"tasmotamanager/utils"
	"time"
)

const (
	defaultHost               = "localhost"
	defaultPort               = "1883"
	defaultClientId           = "tasmotaminder"
	defaultRuleConfigLocation = "/etc/tasmotaminder/rules.yaml"
	connectionTimeout         = 5 * time.Second
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
	clientOptions.SetConnectTimeout(connectionTimeout)
	clientOptions.SetPingTimeout(connectionTimeout)

	return clientOptions
}

func getRuleConfigLocation() string {
	return utils.GetEnvOrDefault("RULE_CONFIG_YAML", defaultRuleConfigLocation)
}

func getRuleConfig() types.PlugRules {
	data, readErr := os.ReadFile(getRuleConfigLocation())
	if readErr != nil {
		log.Fatalf("error opening rule config file: %s", readErr)
	}

	var rules types.PlugRules
	if unmarshalErr := yaml.Unmarshal(data, &rules); unmarshalErr != nil {
		log.Fatalf("error unmarshalling rules: %s", readErr)
	}

	return rules
}
