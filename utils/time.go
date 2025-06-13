package utils

import (
	"time"
)

func FormatTimePointer(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
