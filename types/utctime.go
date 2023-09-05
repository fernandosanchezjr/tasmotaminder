package types

import (
	"fmt"
	"strings"
	"time"
)

type UTCTime struct {
	time.Time
}

func (utc *UTCTime) UnmarshalJSON(b []byte) error {
	rawStr := strings.ReplaceAll(string(b), `"`, "")
	parserStr := fmt.Sprintf("%sZ", rawStr)

	if parsed, parseErr := time.Parse(time.RFC3339, parserStr); parseErr != nil {
		return parseErr
	} else {
		utc.Time = parsed
	}

	return nil
}

func (utc *UTCTime) String() string {
	return utc.Format(time.RFC3339)
}
