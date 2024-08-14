package utils

import (
	"fmt"
	"time"
)

func ParseTimestamp(timestamp string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return "", fmt.Errorf("error parsing time: %s", err)
	}
	return parsedTime.Format("02 Jan 2006 15:04"), nil
}
