package manager

import (
	"context"
	"os"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/grafana"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/network"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/prometheus"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/server"
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
	// Represents Docker SDK client used infrastructure management.
	dockerClient *client.Client

	// Represents metrics monitoring server.
	server *server.MonitoringServer

	// Represents network monitoring component manager.
	networkComponent monitoring.MonitoringComponent

	// Represents Grafana monitoring component manager.
	grafanaComponent monitoring.MonitoringComponent

	// Represents Prometheus monitoring component manager.
	prometheusComponent monitoring.MonitoringComponent
}

// Deploy performs a deployment of monitoring infrastructure.
func (mm *MonitoringManager) Deploy(callback func(err error)) {
	go func() {
		_, err := mm.dockerClient.Ping(context.Background())
		if err != nil {
			callback(errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()))

			return
		}

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

		if !mm.networkComponent.IsDeployed() {
			if err := mm.networkComponent.Deploy(); err != nil {
				callback(errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()))

				return
			}
		}

		err = mm.grafanaComponent.Init()
		if err != nil {
			callback(errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()))

			return
		}

		err = mm.grafanaComponent.Deploy()
		if err != nil {
			callback(errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()))

			return
		}

		err = mm.prometheusComponent.Init()
		if err != nil {
			callback(errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()))

			return
		}

		err = mm.prometheusComponent.Deploy()
		if err != nil {
			if err := mm.grafanaComponent.Remove(); err != nil {
				callback(errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()))

				return
			}

			if err := mm.networkComponent.Remove(); err != nil {
				callback(errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()))

				return
			}

			callback(errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()))

			return
		}

		mm.server.Start(func() {
			if err := mm.grafanaComponent.Remove(); err != nil {
				logging.GetInstance().Fatal(
					errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()).Error())

				return
			}

			if err := mm.prometheusComponent.Remove(); err != nil {
				logging.GetInstance().Fatal(
					errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()).Error())

				return
			}

			if err := mm.networkComponent.Remove(); err != nil {
				logging.GetInstance().Fatal(
					errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()).Error())

				return
			}

			logging.GetInstance().Fatal(ErrMonitoringDeploymentFailed.Error())
		})

		callback(nil)
	}()
}

// Remove performs a removal of monitoring infrastructure.
func (mm *MonitoringManager) Remove(callback func(err error)) {
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

		if mm.networkComponent.IsDeployed() {
			if err := mm.networkComponent.Remove(); err != nil {
				callback(errors.Wrap(err, ErrMonitoringDeploymentFailed.Error()))

				return
			}
		}

		callback(nil)
	}()
}

// newMonitoringManager initializes MonitoringManager.
func newMonitoringManager() *MonitoringManager {
	opts := []client.Opt{
		client.WithAPIVersionNegotiation(),
	}

	if os.Getenv("DOCKER_HOST") != "" {
		opts = append(opts, client.FromEnv)
	} else {
		opts = append(opts,
			client.WithHost("unix://"+os.Getenv("HOME")+"/.docker/run/docker.sock"),
		)
	}

	dockerClient, err := client.NewClientWithOpts(opts...)
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	return &MonitoringManager{
		dockerClient:        dockerClient,
		server:              server.NewMonitoringServer(),
		networkComponent:    network.NewNetworkComponent(dockerClient),
		grafanaComponent:    grafana.NewGrafanaComponent(dockerClient),
		prometheusComponent: prometheus.NewPrometheusComponent(dockerClient),
	}
}
