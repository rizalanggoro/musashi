package utils

import (
	"strconv"
)

// ParseID parses an ID string into a uint.
func ParseID(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

// FormatID formats a uint into an ID string.
func FormatID(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

