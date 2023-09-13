package utils

import (
	"log"
	"strings"
)

func splitTopic(topic string) []string {
	parts := strings.Split(topic, "/")
	if len(parts) != 3 {
		log.Fatalf("invalid topic: %s", topic)
	}
	return parts
}

func GetDeviceId(topic string) string {
	parts := splitTopic(topic)

	return parts[1]
}

func GetTopicSuffix(topic string) string {
	parts := splitTopic(topic)

	return parts[2]
}
