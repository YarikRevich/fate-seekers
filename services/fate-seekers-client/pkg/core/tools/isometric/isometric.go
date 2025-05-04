package isometric

import (
	"golang.org/x/image/math/f64"
)

// ConvertCartesianToIsometric converts provided position according to object size in isometric manner.
func ConvertCartesianToIsometric(xRaw, yRaw int, size float64) f64.Vec2 {
	x := float64(xRaw)
	y := float64(yRaw)

	return f64.Vec2{((x - y) * (size / 2)), (x + y) * (size / 4)}
}

// ConvertIsometricToCartesian converts provided position according to object size in cartesian manner.
func ConvertIsometricToCartesian(raw f64.Vec2, size float64) (xRaw, yRaw int) {
	return int(raw[0]/(size/2)+raw[1]/(size/4)) / 2, int(raw[1]/(size/4)-(raw[0]/(size/2))) / 2
}
