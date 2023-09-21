package utils

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

const (
	tokenWaitTimeout = time.Second * 30
)

func WaitForToken(token mqtt.Token) {
	for !token.WaitTimeout(tokenWaitTimeout) {
		if err := token.Error(); err != nil {
			log.Fatalf("error waiting for token: %s", err)
		}
	}
	if err := token.Error(); err != nil {
		log.Fatalf("error after waiting for token: %s", err)
	}
}
