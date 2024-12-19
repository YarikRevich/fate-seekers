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
