package grafana

import (
	"context"
	"fmt"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/template"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
)

var (
	ErrGrafanaDeployment = errors.New("err happened during grafana deployment")
)

// GrafanaComponent represents a Grafana monitoring component.
type GrafanaComponent struct {
	// Represents Docker SDK client used infrastructure management.
	dockerClient *client.Client
}

func (gc *GrafanaComponent) Init() error {
	return template.Process(
		config.GetDiagnosticsGrafanaConfigDatasourcesDirectory(),
		config.GRAFANA_CONFIG_DATASOURCES_DIAGNOSTICS_TEMPLATE,
		config.GRAFANA_CONFIG_DATASOURCES_DIAGNOSTICS_OUTPUT,
		map[string]interface{}{
			"prometheus": map[string]interface{}{
				"host": config.GetSettingsMonitoringPrometheusName(),
				"port": config.PROMETHEUS_PORT,
			},
		})
}

func (gc *GrafanaComponent) Deploy() error {
	ctx := context.Background()

	_, err := gc.dockerClient.ImagePull(ctx, config.GRAFANA_IMAGE, image.PullOptions{})
	if err != nil {
		return errors.Wrap(err, ErrGrafanaDeployment.Error())
	}

	c := &container.Config{
		Image: config.GRAFANA_IMAGE,
		ExposedPorts: nat.PortSet{
			nat.Port(config.GRAFANA_PORT): struct{}{},
		},
		Env: []string{
			fmt.Sprintf("GF_SECURITY_ADMIN_USER=%s", config.GetSettingsMonitoringGrafanaAdminLogin()),
			fmt.Sprintf("GF_SECURITY_ADMIN_PASSWORD=%s", config.GetSettingsMonitoringGrafanaAdminPassword()),
			"GF_USERS_ALLOW_SIGN_UP=false",
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(config.GRAFANA_PORT): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: config.GRAFANA_PORT,
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   config.GetDiagnosticsGrafanaConfigDirectory(),
				Target:   "/etc/grafana/provisioning/",
				ReadOnly: true,
			},
			{
				Type:   mount.TypeBind,
				Source: config.GetDiagnosticsGrafanaInternalDirectory(),
				Target: "/var/lib/grafana",
			},
		},
	}

	networkingConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			config.GetSettingsMonitoringNetworkName(): {},
		},
	}

	resp, err := gc.dockerClient.ContainerCreate(ctx, c, hostConfig, networkingConfig, nil, config.GetSettingsMonitoringGrafanaName())
	if err != nil {
		return errors.Wrap(err, ErrGrafanaDeployment.Error())
	}

	if err := gc.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return errors.Wrap(err, ErrGrafanaDeployment.Error())
	}

	return nil
}

func (gc *GrafanaComponent) IsDeployed() bool {
	_, err := gc.dockerClient.ContainerInspect(context.Background(), config.GetSettingsMonitoringGrafanaName())

	return err == nil
}

func (gc *GrafanaComponent) Remove() error {
	err := gc.dockerClient.ContainerRemove(
		context.Background(), config.GetSettingsMonitoringGrafanaName(), container.RemoveOptions{
			Force: true,
		})

	if client.IsErrNotFound(err) {
		return nil
	}

	return err
}

// NewGrafanaComponent initializes GrafanaComponent.
func NewGrafanaComponent(dockerClient *client.Client) monitoring.MonitoringComponent {
	return &GrafanaComponent{
		dockerClient: dockerClient,
	}
}
