package sessionname

import "regexp"

// Describes allowed session name validation value.
const (
	allowedPattern = `^[a-zA-Z0-9-]{8}$`
)

// Validate performs provided session name value validation.
func Validate(value string) bool {
	match, err := regexp.MatchString(allowedPattern, value)
	if err != nil {
		return false
	}

	return match
}
