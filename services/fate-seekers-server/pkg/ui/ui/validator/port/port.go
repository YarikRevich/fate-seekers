package port

import "regexp"

// Describes allowed port validation value.
const (
	allowedPattern = `^(?P<port>\d{1,5})$`
)

// Validate performs provided port value validation.
func Validate(value string) bool {
	match, err := regexp.MatchString(allowedPattern, value)
	if err != nil {
		return false
	}

	return match
}
