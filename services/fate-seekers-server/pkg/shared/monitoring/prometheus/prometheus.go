package prometheus

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring"
	"github.com/docker/docker/client"
)

// PrometheusComponent represents a Prometheus monitoring component.
type PrometheusComponent struct {
	// Represents Docker SDK client used infrastructure management.
	dockerClient *client.Client
}

func (pc *PrometheusComponent) Deploy() error {
	return nil
}

func (pc *PrometheusComponent) IsDeployed() bool {
	return false
}

func (pc *PrometheusComponent) Remove() error {
	return nil
}

// NewPrometheusComponent initializes PrometheusComponent.
func NewPrometheusComponent(dockerClient *client.Client) monitoring.MonitoringComponent {
	return &PrometheusComponent{
		dockerClient: dockerClient,
	}
}
