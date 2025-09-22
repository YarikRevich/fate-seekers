package common

import "errors"

// Describes all the common errors for sound management system.
var (
	ErrSoundPlayerAccess = errors.New("err happened during sound player access")
)

const (
	// Represents player bytes per sample value.
	BytesPerSample = 8

	// Represents player sample rate value.
	SampleRate = 44100
)
