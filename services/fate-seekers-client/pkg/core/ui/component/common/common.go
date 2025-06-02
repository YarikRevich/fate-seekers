package common

import (
	"fmt"
	"image/color"
)

var (
	// Represents common button text color.
	ButtonTextColor = color.RGBA{R: 11, G: 16, B: 37, A: 255}

	// Represents common notification error text color.
	NotificationErrorTextColor = color.RGBA{R: 245, G: 0, B: 0, A: 255}

	// Represents common notification info text color.
	NotificationInfoTextColor = color.White
)

// ComposeMessage performs message composition.
func ComposeMessage(err string, details string) string {
	return fmt.Sprintf("%s\n Details: %s", err, details)
}
