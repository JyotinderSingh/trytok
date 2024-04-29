package docker

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	imageName      = "jyotindersingh/compiler-server"
	hostIP         = "0.0.0.0"
	hostPort       = "5500"
	containerPort  = "8080/tcp"
	httpTimeout    = 10 * time.Second
	checkInterval  = 500 * time.Millisecond
	startupTimeout = 30 * time.Second
)

// CreateContainer starts a new container and waits until it's reachable via HTTP.
func CreateContainer() (string, *client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", nil, err
	}

	containerID, err := createAndStartContainer(cli)
	if err != nil {
		return "", nil, err
	}

	if err := waitForContainerReady(); err != nil {
		return "", nil, err
	}

	return containerID, cli, nil
}

func createAndStartContainer(cli *client.Client) (string, error) {
	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{"go", "run", "/server.go"},
		Tty:   false,
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				{HostIP: hostIP, HostPort: hostPort},
			},
		},
	}, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

func waitForContainerReady() error {
	timeout := time.After(startupTimeout)
	tick := time.NewTicker(checkInterval)
	defer tick.Stop()

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s", hostIP, hostPort), bytes.NewBufferString(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	httpClient := &http.Client{Timeout: httpTimeout}

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for container to be ready")
		case <-tick.C:
			res, err := httpClient.Do(req)
			if err == nil && res.StatusCode == 200 {
				res.Body.Close()
				return nil
			}
			if res != nil {
				res.Body.Close()
			}
		}
	}
}
