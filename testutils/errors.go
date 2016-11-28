package testutils

import (
	"regexp"
)

// IsError returns true if err is non-nil and the error string matches the
// supplied regexp.
func IsError(err error, re string) bool {
	if err == nil {
		return false
	}
	matched, merr := regexp.MatchString(re, err.Error())
	if merr != nil {
		return false
	}
	return matched
}
