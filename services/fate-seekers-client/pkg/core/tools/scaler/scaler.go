package scaler

// GetScaleFactor retrieveds scale factor for the given sizes.
func GetScaleFactor(currentSize int, expectedSize int) float64 {
	return float64(expectedSize) / float64(currentSize)
}
