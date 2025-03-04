package host

import "regexp"

// Describes allowed host validation value.
const (
	allowedPattern = `^(?P<host>([a-zA-Z0-9]+\.)*[a-zA-Z0-9]+\.[a-zA-Z]{2,63}|localhost|\b\d{1,3}(\.\d{1,3}){3}):(?P<port>\d{1,5})$`
)

// Validate performs provided host value validation.
func Validate(value string) bool {
	match, err := regexp.MatchString(allowedPattern, value)
	if err != nil {
		return false
	}

	return match
}
