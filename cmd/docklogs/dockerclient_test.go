package docklogs

import (
	"context"
	"github.com/docker/docker/api/types"
	"os"
	"testing"
)

type containerListFn func(types.ContainerListOptions) ([]types.Container, error)
type fakeClient struct {
	*DockerClient
	containerListFunc containerListFn
}

func NewFakeClient() *fakeClient {
	return &fakeClient{
		&DockerClient{},
		nil,
	}
}

func (f *fakeClient) WithContainers(names ...string) *fakeClient {
	for _, name := range names {
		c := types.Container{
			Names: []string{"/" + name},
		}
		f.LiveContainers = append(f.LiveContainers, c)
	}

	return f
}

func (f *fakeClient) WithContainerListFunc(fn containerListFn) *fakeClient {
	f.containerListFunc = fn
	return f
}

func (f *fakeClient) ContainerList(_ context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	return f.LiveContainers, nil
}

func TestDockerClient_IsContainerLive(t *testing.T) {
	cli := NewFakeClient().WithContainers("a", "b")
	tests := []struct {
		name          string
		containerName string
		want          bool
	}{
		{
			name:          "live container",
			containerName: "a",
			want:          true,
		},
		{
			name:          "dead container",
			containerName: "c",
			want:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cli.IsContainerLive(tt.containerName); got != tt.want {
				t.Errorf("IsContainerLive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDockerClient_Successful(t *testing.T) {
	got, err := GetDockerClient(context.TODO())
	if err != nil {
		t.Errorf("expected nil got err: %v", err)
	}

	if got == nil {
		t.Errorf("client appears should not be nil")
	}
}

func TestGetDockerClient_DockerNotRunning(t *testing.T) {
	os.Setenv("_FORCE_DOCKER_RUNNING", "False")
	defer os.Unsetenv("_FORCE_DOCKER_RUNNING")

	got, err := GetDockerClient(context.TODO())
	if err == nil {
		t.Errorf("expected error got nil")
	}

	if got != nil {
		t.Errorf("expected nil docker client got: %v", got)
	}
}
