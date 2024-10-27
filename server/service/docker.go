package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerService struct {
	cli *client.Client
}

func NewDockerService() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerService{cli: cli}, nil
}

func (ds *DockerService) RestartContainer(containerID string) error {
	ctx := context.Background()
	opts := container.StopOptions{}
	if err := ds.cli.ContainerRestart(ctx, containerID, opts); err != nil {
		return fmt.Errorf("failed to restart container %s: %w", containerID, err)
	}

	return nil
}

func (ds *DockerService) GetContainerInfo(containerID string) (*types.ContainerJSON, error) {
	ctx := context.Background()

	containerJSON, err := ds.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container %s: %w", containerID, err)
	}

	return &containerJSON, nil
}

func (ds *DockerService) ListContainers() ([]types.Container, error) {
	ctx := context.Background()
	opts := container.ListOptions{}
	containers, err := ds.cli.ContainerList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	return containers, nil
}
