package docker

import (
	"bytes"
	"context"
	"log"
	"time"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type Manager struct {
	cli *client.Client
}

func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	_, err = cli.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to Docker daemon")
	return &Manager{cli: cli}, nil
}

func (m *Manager) RunCode (code string) (string , error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	containerConfig := &container.Config{
		Image: "python:3.12-alpine",
		Cmd: []string{"python", "-c", code},
		Tty: false,
	}

	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			Memory: 100*1024*1024,
		},
		NetworkMode: "none",
		ReadonlyRootfs: true,
		AutoRemove: true,
	}

	resp, err := m.cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err!=nil {
		return "", err
	}
	log.Printf("Container created with ID: %s", resp.ID)

	if err := m.cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		return "", err
	}
	log.Printf("Container %s started", resp.ID)

	statusCh, errorCh := m.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
		case err := <-errorCh: {
			if err!=nil {
				return "", err
			}
		}
		case <-statusCh: {
			log.Printf("Container %s finished", resp.ID)
		}
	}

	out, err := m.cli.ContainerLogs(ctx, resp.ID, client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", err
	}
	defer out.Close()

	var stdOut, stdErr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdOut, &stdErr, out)
	if err != nil {
		return "", err
	}

	if stdErr.Len() > 0 {
        return stdErr.String(), nil
    }

	return stdOut.String(), nil

}