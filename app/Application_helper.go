package app

import (
	"fmt"
	"time"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
}

