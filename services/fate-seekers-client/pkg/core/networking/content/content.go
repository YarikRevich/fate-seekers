package content

// NetworkingContentChannel represent networking content channel interface.
type NetworkingContentChannel interface {
	// Schedule performs channel call once.
	Schedule(args interface{}, finishCallback func())
}
