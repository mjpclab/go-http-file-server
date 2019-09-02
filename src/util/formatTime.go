package util

import (
	"time"
)

func FormatTimeMinute(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func FormatTimeSecond(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
