package monitoring

// MonitoringComponent represents monitoring component interface.
type MonitoringComponent interface {
	// Init performs component templates initialization.
	Init() error

	// Deploy performs monitoring
	Deploy() error

	// IsDeployed checks if monitoring component is deployed.
	IsDeployed() bool

	// Remove performs monitoring component removal.
	Remove() error
}
