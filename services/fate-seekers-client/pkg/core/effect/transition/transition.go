package transition

// TransitionEffect represents transition effects interface.
type TransitionEffect interface {
	// Done checks if transition has been finished.
	Done() bool

	// OnEnd checks if transition is on end state.
	OnEnd() bool

	// Update handles transition state update.
	Update()

	// Clean performes forced memory cleanup for the transition only.
	Clean()

	// Reset performes transition state reset.
	Reset()

	// GetValue retrieves updated value.
	GetValue() float64
}
