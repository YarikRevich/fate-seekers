package networking

// NetworkingChannel represent channel interface.
type NetworkingChannel interface {
	// Execute performs channel call.
	Execute(args interface{}, finishCallback func())
}
