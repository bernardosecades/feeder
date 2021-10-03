package numberfn

import (
	"strconv"
)

// IsUnsignedInteger check if string is a unsigned integer
func IsUnsignedInteger(text string) bool {
	v, err := strconv.Atoi(text)
	if err != nil {
		return false
	}

	if v < 0 {
		return false
	}

	return true
}
