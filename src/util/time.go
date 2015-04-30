package util

import (
	"time"
)

func FormalTime(timestamp int64) string {
	const layout = "2006-01-02 15:04:05"
	timeChange := time.Unix(timestamp, 0)
	return timeChange.Format(layout)
}

func NowTime() string {
	const layout = "2006-01-02 15:04:05"
	return time.Now().Format(layout)
}
