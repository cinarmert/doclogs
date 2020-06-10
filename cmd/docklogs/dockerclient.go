package docklogs

import (
	"context"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"os"
)

// DockerCli extends the docker.APIClient with IsContainerLive function
// in order to reduce the time to check if a container is live.
type DockerCli interface {
	docker.APIClient
	IsContainerLive(string) bool
}

// DockerClient extends the docker.Client with LiveContainers field, in order to
// reduce the time to check if a container is live.
type DockerClient struct {
	*docker.Client
	LiveContainers []types.Container // Used to check if the given containers are live
}

// GetDockerClient returns a pointer to new DockerClient, it fetches all the live containers
// on the creation of the client in order to reduce the time to check if a container is live.
func GetDockerClient(ctx context.Context) (*DockerClient, error) {
	if os.Getenv("_FORCE_DOCKER_RUNNING") == "False" { // for testing
		return nil, errors.New("docker daemon not running")
	}

	dockerClient, err := docker.NewClientWithOpts(docker.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrap(err, "could not create docker client, make sure the docker daemon is running")
	}

	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "could not get list of containers")
	}

	return &DockerClient{
		dockerClient,
		containers,
	}, nil
}

// IsContainerLive checks if the container with given name is alive or dead
func (dc *DockerClient) IsContainerLive(name string) bool {
	for _, c := range dc.LiveContainers {
		for _, n := range c.Names {
			if "/"+name == n {
				return true
			}
		}
	}

	return false
}
