package external

// Represents external cursor position source.
var externalCursorPosition func() (int, int)

// SetExternalCursorPositionSource sets external cursor position source.
func SetExternalCursorPositionSource(value func() (int, int)) {
	externalCursorPosition = value
}

// GetExternalCursorPositionSource retrieves external cursor position source.
func GetExternalCursorPositionSource() func() (int, int) {
	return externalCursorPosition
}

// Represents external left mouse click source.
var externalLeftMouseClick func() bool

// SetExternalLeftMouseClick sets external left mouse click source.
func SetExternalLeftMouseClick(value func() bool) {
	externalLeftMouseClick = value
}

// GetExternalLeftMouseClick retrieves external left mouse click source.
func GetExternalLeftMouseClick() func() bool {
	return externalLeftMouseClick
}

// Represents external right mouse click source.
var externalRightMouseClick func() bool

// SetExternalRightMouseClick sets external right mouse click source.
func SetExternalRightMouseClick(value func() bool) {
	externalRightMouseClick = value
}

// GetExternalRightMouseClick retrieves external left mouse click source.
func GetExternalRightMouseClick() func() bool {
	return externalRightMouseClick
}

// Represents external middle mouse click source.
var externalMiddleMouseClick func() bool

// SetExternalMiddleMouseClick sets external middle mouse click source.
func SetExternalMiddleMouseClick(value func() bool) {
	externalMiddleMouseClick = value
}

// GetExternalMiddleMouseClick retrieves external middle mouse click source.
func GetExternalMiddleMouseClick() func() bool {
	return externalMiddleMouseClick
}

// Represents external wheel mouse source.
var externalWheelMouse func() (float64, float64)

// SetExternalWheelMouse sets external wheel mouse source.
func SetExternalWheelMouse(value func() (float64, float64)) {
	externalWheelMouse = value
}

// GetExternalWheelMouse retrieves external wheel mouse source.
func GetExternalWheelMouse() func() (float64, float64) {
	return externalWheelMouse
}
