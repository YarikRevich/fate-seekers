package prometheus

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/template"
	"github.com/docker/docker/client"
)

// PrometheusComponent represents a Prometheus monitoring component.
type PrometheusComponent struct {
	// Represents Docker SDK client used infrastructure management.
	dockerClient *client.Client
}

func (pc *PrometheusComponent) Init() error {
	return template.Process(
		config.GetDiagnosticsPrometheusConfigDirectory(),
		config.PROMETHEUS_CONFIG_DIAGNOSTICS_TEMPLATE,
		config.PROMETHEUS_CONFIG_DIAGNOSTICS_OUTPUT,
		map[string]interface{}{
			"metrics": map[string]interface{}{
				"host": "host.docker.internal",
				"port": config.GetSettingsMonitoringPrometheusPort(),
			},
		})
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
