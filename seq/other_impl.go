package seq

import (
	"strings"
)

// overlap2 is another implementation of the overlap function.
// This function uses strings.Compare.
func overlap2(x, y string) int {
	maxLen := len(x)
	if len(y) < maxLen {
		maxLen = len(y)
	}
	minLen := len(x)/2 + 1
	if l := len(y)/2 + 1; l > minLen {
		minLen = l
	}

	for n := maxLen - 1; n >= minLen; n-- {
		if strings.Compare(x[len(x)-n:], y[:n]) == 0 {
			return n
		}
	}
	return -1
}
