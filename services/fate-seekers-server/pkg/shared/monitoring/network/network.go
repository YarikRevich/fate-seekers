package network

import (
	"context"
	"errors"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var (
	ErrGrafanaDeployment = errors.New("err happened during grafana deployment")
)

// NetworkComponent represents a networking monitoring component.
type NetworkComponent struct {
	// Represents Docker SDK client used infrastructure management.
	dockerClient *client.Client
}

func (nc *NetworkComponent) Init() error {
	return nil
}

func (nc *NetworkComponent) Deploy() error {
	ctx := context.Background()

	ipamConfig := []network.IPAMConfig{
		{
			Subnet: "149.156.139.0/28",
		},
	}

	createOpts := network.CreateOptions{
		Driver: "bridge",
		IPAM: &network.IPAM{
			Config: ipamConfig,
		},
	}

	_, err := nc.dockerClient.NetworkCreate(
		ctx, config.GetSettingsMonitoringNetworkName(), createOpts)

	return err
}

func (nc *NetworkComponent) IsDeployed() bool {
	_, err := nc.dockerClient.NetworkInspect(
		context.Background(), config.GetSettingsMonitoringNetworkName(), network.InspectOptions{})

	return err == nil
}

func (nc *NetworkComponent) Remove() error {
	return nc.dockerClient.NetworkRemove(
		context.Background(), config.GetSettingsMonitoringNetworkName())
}

// NewNetworkComponent initializes NetworkComponent.
func NewNetworkComponent(dockerClient *client.Client) monitoring.MonitoringComponent {
	return &NetworkComponent{
		dockerClient: dockerClient,
	}
}
