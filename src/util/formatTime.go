package util

import "time"

const timeLayout = "2006-01-02 15:04:05"

func FormatTimeSecond(t time.Time) string {
	return t.Format(timeLayout)
}

func AppendTimeSecond(buf []byte, t time.Time) []byte {
	return t.AppendFormat(buf, timeLayout)
}
