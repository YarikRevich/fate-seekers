package manager

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/grafana"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/prometheus"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

var (
	// GetInstance retrieves instance of the networking content connector, performing initilization if needed.
	GetInstance = sync.OnceValue[*MonitoringManager](newMonitoringManager)
)

var (
	ErrMonitoringDeploymentFailed = errors.New("err happened during monitoring deployment")
)

// MonitoringManager represents monitoring manager.
type MonitoringManager struct {
	// Represents Grafana monitoring component manager.
	grafanaComponent monitoring.MonitoringComponent

	// Represents Prometheus monitoring component manager.
	prometheusComponent monitoring.MonitoringComponent
}

// Deploy performs a deployment of monitoring infrastructure.
func (mm *MonitoringManager) Deploy(callback func(err error)) {
	go func() {
		if mm.grafanaComponent.IsDeployed() {
			if err := mm.grafanaComponent.Remove(); err != nil {
				callback(errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()))

				return
			}
		}

		if mm.prometheusComponent.IsDeployed() {
			if err := mm.prometheusComponent.Remove(); err != nil {
				callback(errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()))

				return
			}
		}

		mm.grafanaComponent.Deploy()
		// f, err := dockerClient.ContainerList(context.Background(), container.ListOptions{})
		// fmt.Println(f, err)

		// f[0].Names

		// if mm.grafanaComponent.Deploy() {

		// }

		callback(nil)
	}()
}

// Remove performs a removal of monitoring infrastructure.
func (mm *MonitoringManager) Remove(callback func(err error)) {
	go func() {
		callback(nil)
	}()
}

// newMonitoringManager initializes MonitoringManager.
func newMonitoringManager() *MonitoringManager {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	return &MonitoringManager{
		grafanaComponent:    grafana.NewGrafanaComponent(dockerClient),
		prometheusComponent: prometheus.NewPrometheusComponent(dockerClient),
	}
}
