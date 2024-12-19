package dto

// ReducerResult represents result of reducer execution operation.
type ReducerResult map[string]interface{}

// ReducerResultUnit represents result unit of reducer execution operation.
type ReducerResultUnit struct {
	// Represents reducer result key.
	Key string

	// Represents reducer result value.
	Value interface{}
}

// ComposeReducerResult composes reducer result from the given reducer result units.
func ComposeReducerResult(units ...ReducerResultUnit) ReducerResult {
	result := make(map[string]interface{})

	for _, unit := range units {
		result[unit.Key] = unit.Value
	}

	return result
}
