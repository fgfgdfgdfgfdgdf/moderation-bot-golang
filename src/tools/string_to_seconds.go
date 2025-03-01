package tools

import (
	"strconv"
)

func StringToSeconds(duration string) int64 {
	var totalSeconds int64

	number, _ := strconv.ParseInt(duration[:len(duration)-1], 10, 64)

	unit := duration[len(duration)-1]

	switch unit {
	case 's':
		totalSeconds = number
	case 'm':
		totalSeconds = number * 60
	case 'h':
		totalSeconds = number * 3600
	case 'd':
		totalSeconds = number * 86400
	case 'w':
		totalSeconds = number * 604800
	}

	return totalSeconds
}
