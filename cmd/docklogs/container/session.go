package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"io"
	"sync"
)

// Session holds the information for a live container,
// a socket will be opened through a session to fetch the logs.
type Session struct {
	Name         string
	Status       Status
	dockerClient docker.APIClient
	follow       bool
	ctx          context.Context
}

// NewContainerSession creates a new session for a given container name.
func NewContainerSession(dockerClient docker.APIClient, ctx context.Context, name string, follow bool) *Session {
	return &Session{
		Name:         name,
		dockerClient: dockerClient,
		ctx:          ctx,
		follow:       follow,
	}
}

// ReadLogs fetches the logs and writes it into provided writer w,
// it follows the logs if the follow option is enabled by the user.
func (c *Session) ReadLogs(wg *sync.WaitGroup, w io.Writer) error {
	rc, err := c.dockerClient.ContainerLogs(c.ctx, c.Name, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: false,
		Follow:     c.follow,
		Tail:       "all",
	})
	c.Status = Live

	if err != nil {
		c.Status = Errored
		return errors.Wrap(err, fmt.Sprintf("could not read logs from container: %v", err))
	}

	_, err = io.Copy(w, rc)
	if err != nil {
		c.Status = Errored
		return errors.Wrap(err, fmt.Sprintf("could not copy from container socket: %v", err))
	}

	c.Status = Terminated
	wg.Done()
	return nil
}
