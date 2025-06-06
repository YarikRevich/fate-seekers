package sessionid

import "regexp"

// Describes allowed sessionid validation value.
const (
	allowedPattern = `^[a-zA-Z0-9-]{8}\b$`
)

// Validate performs provided sessionid value validation.
func Validate(value string) bool {
	match, err := regexp.MatchString(allowedPattern, value)
	if err != nil {
		return false
	}

	return match
}

// rand.Seed(time.Now().UnixNano())
// 	return fmt.Sprintf("%08d", rand.Intn(100000000))
