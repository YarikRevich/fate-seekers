package grafana

import (
	"context"
	"fmt"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring"
	"github.com/docker/cli/cli/command/image"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	grafanaContainerName = "fate-seekers-grafana"
	grafanaImage         = "grafana/grafana:latest"
	grafanaPort          = "3000"
)

// GrafanaComponent represents a Grafana monitoring component.
type GrafanaComponent struct {
	// Represents Docker SDK client used infrastructure management.
	dockerClient *client.Client
}

func (gc *GrafanaComponent) Deploy() error {
	ctx := context.Background()

	// 1. Pull the image
	_, err := gc.dockerClient.ImagePull(ctx, grafanaImage, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull grafana image: %w", err)
	}

	config := &container.Config{
		Image: grafanaImage,
		ExposedPorts: nat.PortSet{
			nat.Port(grafanaPort): struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(grafanaPort): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: grafanaPort,
				},
			},
		},
		AutoRemove: true,
	}

	resp, err := gc.dockerClient.ContainerCreate(ctx, config, hostConfig, nil, nil, grafanaContainerName)
	if err != nil {
		return fmt.Errorf("failed to create grafana container: %w", err)
	}

	if err := gc.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start grafana container: %w", err)
	}

	return nil
}

func (gc *GrafanaComponent) IsDeployed() bool {
	_, err := gc.dockerClient.ContainerInspect(context.Background(), grafanaContainerName)

	return err == nil
}

func (gc *GrafanaComponent) Remove() error {
	ctx := context.Background()

	// Force remove the container (stops it if running)
	err := gc.dockerClient.ContainerRemove(ctx, grafanaContainerName, container.RemoveOptions{
		Force: true,
	})

	// If the error is simply that the container doesn't exist, we can ignore it
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
