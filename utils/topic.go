package utils

import (
	"log"
	"strings"
)

func GetDeviceId(topic string) string {
	parts := strings.Split(topic, "/")
	if len(parts) != 3 {
		log.Fatalf("invalid topic: %s", topic)
	}

	return parts[1]
}
