package encryptionkey

// Validate performs provided encryption key value validation.
func Validate(value string) bool {
	return len(value) <= 64
}
