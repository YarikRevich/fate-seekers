package scaler

// GetScaleFactor retrieves scale factor for the given sizes.
func GetScaleFactor(currentSize int, expectedSize int) float64 {
	return float64(expectedSize) / float64(currentSize)
}

// GetPercentageOf retrieves given percentage of the given size.
func GetPercentageOf(size int, percentage int) int {
	return (size * percentage) / 100
}
