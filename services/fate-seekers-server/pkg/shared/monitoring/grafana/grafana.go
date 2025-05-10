package grafana

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring"
	"github.com/docker/docker/client"
)

// GrafanaComponent represents a Grafana monitoring component.
type GrafanaComponent struct {
	// Represents Docker SDK client used infrastructure management.
	dockerClient *client.Client
}

func (gc *GrafanaComponent) Deploy() error {
	return nil
}

func (gc *GrafanaComponent) IsDeployed() bool {
	return false
}

func (gc *GrafanaComponent) Remove() error {
	return nil
}

// NewGrafanaComponent initializes GrafanaComponent.
func NewGrafanaComponent(dockerClient *client.Client) monitoring.MonitoringComponent {
	return &GrafanaComponent{
		dockerClient: dockerClient,
	}
}
