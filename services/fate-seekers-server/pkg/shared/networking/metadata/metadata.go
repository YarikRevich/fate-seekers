package metadata

// NetworkingMetadataChannel represent networking metadata channel interface.
type NetworkingMetadataChannel interface {
	// ScheduleOnce performs channel call once.
	ScheduleOnce(args interface{}, finishCallback func())
}
