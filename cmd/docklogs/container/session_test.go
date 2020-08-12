package container

import (
	"bytes"
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"gotest.tools/assert"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"testing"
)

type logfn func(string, types.ContainerLogsOptions) (io.ReadCloser, error)
type fakeClient struct {
	docker.Client
	logFunc logfn
}

func (f *fakeClient) ContainerLogs(_ context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	if f.logFunc != nil {
		return f.logFunc(container, options)
	}
	return nil, nil
}

func TestNewContainerSession(t *testing.T) {
	type args struct {
		dockerClient *docker.Client
		ctx          context.Context
		name         string
		follow       bool
	}
	tests := []struct {
		name string
		args args
		want *Session
	}{
		{
			name: "nil client, nil context",
			args: args{
				name: "test",
			},
			want: &Session{
				Name: "test",
			},
		},
		{
			name: "nil client, todo ctx",
			args: args{
				ctx:    context.TODO(),
				name:   "test",
				follow: true,
			},
			want: &Session{
				ctx:    context.TODO(),
				Name:   "test",
				follow: true,
			},
		},
		{
			name: "empty client, todo ctx",
			args: args{
				dockerClient: &docker.Client{},
				ctx:          context.TODO(),
				name:         "test",
			},
			want: &Session{
				dockerClient: &docker.Client{},
				ctx:          context.TODO(),
				Name:         "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewContainerSession(tt.args.dockerClient, tt.args.ctx, tt.args.name, tt.args.follow, 100)
			assert.Equal(t, got.Name, tt.want.Name)
			assert.Equal(t, got.follow, tt.want.follow)
		})
	}
}

func TestSession_ReadLogs(t *testing.T) {
	tests := []struct {
		name       string
		logFunc    logfn
		out        *bytes.Buffer
		want       string
		wantErr    bool
		wantStatus Status
	}{
		{
			name: "successful attempt",
			logFunc: func(s string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
				return ioutil.NopCloser(strings.NewReader("test log")), nil
			},
			out:        new(bytes.Buffer),
			want:       "test log",
			wantErr:    false,
			wantStatus: Terminated,
		},
		{
			name: "log errored",
			logFunc: func(s string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
				return nil, errors.New("failed")
			},
			out:        new(bytes.Buffer),
			want:       "",
			wantErr:    true,
			wantStatus: Errored,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fakeClient{
				logFunc: tt.logFunc,
			}

			s := NewContainerSession(f, context.TODO(), "test", false, 100)
			wg := &sync.WaitGroup{}
			wg.Add(1)
			err := s.ReadLogs(wg, tt.out)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.out.String(), tt.want)
			assert.Equal(t, tt.wantStatus, s.Status)
		})
	}
}
