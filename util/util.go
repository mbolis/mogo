package util

import (
	"slices"
)

func NTString(bytes []byte) string {
	i := slices.Index(bytes, 0)
	switch i {
	case -1:
		return string(bytes)
	default:
		return string(bytes[:i])
	}
}
