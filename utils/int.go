package utils

import "time"

func DurationMax(a time.Duration, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
