package port

import (
	"regexp"
	"strconv"
)

// Describes allowed port validation value.
const (
	allowedPattern = `^(?P<port>\d{1,5})$`
)

// Validate performs provided port value validation.
func Validate(value int) bool {
	match, err := regexp.MatchString(allowedPattern, strconv.Itoa(value))
	if err != nil {
		return false
	}

	return match
}
