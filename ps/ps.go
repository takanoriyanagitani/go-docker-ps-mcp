package ps

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func GetContainers(ctx context.Context, socketPath string) ([]types.Container, error) {
	var opts []client.Opt
	if socketPath != "" {
		opts = append(opts, client.WithHost("unix://"+socketPath))
	} else {
		opts = append(opts, client.FromEnv)
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}

	return containers, nil
}
