package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// CreateContainer starts a new container with the provided code snippet.
func CreateContainer(code string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "yourimage",
		Cmd:   []string{"your-compiler-binary", code},
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

// FetchLogs retrieves the logs from a container.
func FetchLogs(containerID string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	out, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", err
	}

	logs, err := io.ReadAll(out)
	if err != nil {
		return "", err
	}

	return string(logs), nil
}
