package sessionseed

import "regexp"

// Describes allowed session seed validation value.
const (
	allowedPattern = `^[0-9-]{8}$`
)

// Validate performs provided session seed value validation.
func Validate(value string) bool {
	match, err := regexp.MatchString(allowedPattern, value)
	if err != nil {
		return false
	}

	return match
}
