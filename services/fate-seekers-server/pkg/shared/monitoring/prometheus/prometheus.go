package prometheus

import (
	"context"
	"fmt"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/template"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
)

var (
	ErrPrometheusDeployment = errors.New("err happened during prometheus deployment")
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
	ctx := context.Background()

	_, err := pc.dockerClient.ImagePull(ctx, config.PROMETHEUS_IMAGE, image.PullOptions{})
	if err != nil {
		return errors.Wrap(err, ErrPrometheusDeployment.Error())
	}

	c := &container.Config{
		Image: config.PROMETHEUS_IMAGE,
		ExposedPorts: nat.PortSet{
			nat.Port(config.PROMETHEUS_PORT): struct{}{},
		},
		Cmd: []string{
			"--config.file=/etc/prometheus/prometheus.yml",
			"--storage.tsdb.path=/prometheus",
			"--web.console.libraries=/usr/share/prometheus/console_libraries",
			"--web.console.templates=/usr/share/prometheus/consoles",
			fmt.Sprintf("--web.listen-address=0.0.0.0:%s", config.PROMETHEUS_PORT),
			"--web.enable-lifecycle",
			"--web.enable-admin-api",
		},
	}

	hostConfig := &container.HostConfig{
		ExtraHosts: []string{
			"host.docker.internal:host-gateway",
		},
		PortBindings: nat.PortMap{
			nat.Port(config.PROMETHEUS_PORT): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: config.PROMETHEUS_PORT,
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   config.GetDiagnosticsPrometheusConfigDirectory(),
				Target:   "/etc/prometheus/",
				ReadOnly: true,
			},
			{
				Type:   mount.TypeBind,
				Source: config.GetDiagnosticsPrometheusInternalDirectory(),
				Target: "/prometheus",
			},
		},
	}

	resp, err := pc.dockerClient.ContainerCreate(ctx, c, hostConfig, nil, nil, config.GetSettingsMonitoringPrometheusName())
	if err != nil {
		return errors.Wrap(err, ErrPrometheusDeployment.Error())
	}

	if err := pc.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return errors.Wrap(err, ErrPrometheusDeployment.Error())
	}

	return nil
}

func (pc *PrometheusComponent) IsDeployed() bool {
	_, err := pc.dockerClient.ContainerInspect(context.Background(), config.GetSettingsMonitoringPrometheusName())

	return err == nil
}

func (pc *PrometheusComponent) Remove() error {
	err := pc.dockerClient.ContainerRemove(
		context.Background(), config.GetSettingsMonitoringPrometheusName(), container.RemoveOptions{
			Force: true,
		})

	if client.IsErrNotFound(err) {
		return nil
	}

	return err
}

// NewPrometheusComponent initializes PrometheusComponent.
func NewPrometheusComponent(dockerClient *client.Client) monitoring.MonitoringComponent {
	return &PrometheusComponent{
		dockerClient: dockerClient,
	}
}
