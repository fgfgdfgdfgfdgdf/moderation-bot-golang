package tools

import (
	"regexp"
	"strconv"
)

var (
	r = regexp.MustCompile(`^[0-9]+(s|m|h|d|w)$`)
)

func Check_duration(s string) bool {

	if !r.MatchString(s) {
		return false
	}

	num, err := strconv.Atoi(s[:len(s)-1])
	if err != nil {
		return false
	}

	return num <= 2000
}
