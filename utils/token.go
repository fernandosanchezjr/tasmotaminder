package utils

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

const (
	DefaultWaitTimeout = time.Second * 5
)

func WaitForToken(token mqtt.Token) {
	for !token.WaitTimeout(DefaultWaitTimeout) {
		if err := token.Error(); err != nil {
			log.Fatalf("error waiting for token: %s", err)
		}
	}
}
