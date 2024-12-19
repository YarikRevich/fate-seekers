package entry

import "sync"

var (
	// GetInstance retrieves instance of the entry screen, performing initilization if needed.
	GetInstance = sync.OnceValue[*EntryScreen](newEntryScreen)
)

// EntryScreen represents entry screen implementation.
type EntryScreen struct {
}

// newEntryScreen initializes EntryScreen.
func newEntryScreen() *EntryScreen {
	return new(EntryScreen)
}
