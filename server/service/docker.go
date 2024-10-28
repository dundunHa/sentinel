package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"sentinel/server/model"
	"strings"
)

type DockerService struct {
	cli *client.Client
}

var DockerClientMap = make(map[model.ProcessorAPPID]*client.Client)

func RegisterDockerClient() error {
	if err := registerClient(1, "unix:///var/run-host/docker.sock"); err != nil {
		return err
	}
	if err := registerClient(4, "http://192.168.1.2:12375"); err != nil {
		return err
	}
	return nil
}

func registerClient(appid model.ProcessorAPPID, addr string) error {
	client, err := NewDockerClient(addr)
	if err != nil {
		return err
	}
	DockerClientMap[appid] = client

	return nil
}

func NewDockerClient(addr string) (*client.Client, error) {
	var cli *client.Client
	var err error

	if strings.HasPrefix(addr, "unix://") {
		// 使用 Unix 套接字方式连接
		cli, err = client.NewClientWithOpts(
			client.WithHost(addr),
			client.WithAPIVersionNegotiation(),
		)
	} else {
		// 使用 TCP 方式连接，确保前缀为 tcp://
		if !strings.HasPrefix(addr, "tcp://") {
			addr = "tcp://" + addr
		}
		cli, err = client.NewClientWithOpts(
			client.WithHost(addr),
			client.WithAPIVersionNegotiation(),
		)
	}

	if err != nil {
		return nil, err
	}

	return cli, nil
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
