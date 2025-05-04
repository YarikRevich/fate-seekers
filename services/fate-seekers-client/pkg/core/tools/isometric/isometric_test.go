package isometric

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestConvertCartesianToIsometric tests cartesian to isometric convertion.
func TestConvertCartesianToIsometric(t *testing.T) {
	xRaw := 10
	yRaw := 10
	size := 20.0

	result := ConvertCartesianToIsometric(xRaw, yRaw, size)

	x := float64(xRaw)
	y := float64(yRaw)

	require.Equal(t, result[0], ((x - y) * (size / 2)))
	require.Equal(t, result[1], (x+y)*(size/4))
}
