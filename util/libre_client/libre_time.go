package libreclient

import (
	"strings"
	"time"
)

const libreTimestampLayout = "1/2/2006 3:04:05 PM"

// Libre time format with 1/2/2006 3:04:05 PM layout.
type LibreTimeFormat struct {
	time.Time
}

func (t *LibreTimeFormat) UnmarshalJSON(data []byte) error {
	timeString := strings.ReplaceAll(string(data), "\"", "")
	parsedTime, err := time.Parse(libreTimestampLayout, timeString)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}
