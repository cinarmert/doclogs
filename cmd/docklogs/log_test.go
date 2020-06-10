package docklogs

import (
	"context"
	"github.com/cinarmert/doclogs/cmd/docklogs/container"
	"github.com/docker/docker/api/types"
	"testing"
)

func TestLogOp_createContainerSessions(t *testing.T) {
	lo := &LogOp{
		Containers: []string{"alpine", "aa"},
	}

	type args struct {
		name    string
		dc      *fakeClient
		wantErr bool
	}

	tests := []args{
		{
			name: "successful",
			dc: NewFakeClient().WithContainers("aa").WithContainerListFunc(func(options types.ContainerListOptions) ([]types.Container, error) {
				return []types.Container{}, nil
			}),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// want graceful shutdown
			lo.createContainerSessions(context.Background(), tt.dc)
		})
	}
}

func TestLogOp_report(t *testing.T) {
	tests := []struct {
		name     string
		sessions []*container.Session
	}{
		{
			name:     "nonempty session list",
			sessions: []*container.Session{{Name: "test", Status: container.Terminated}},
		},
		{
			name:     "empty session list",
			sessions: []*container.Session{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LogOp{}
			l.report(tt.sessions)
		})
	}
}

func TestLogOp_Run(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		Containers []string
		wantErr    bool
	}{
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: true,
		},
		{
			name:    "todo context, empty containers",
			ctx:     context.TODO(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LogOp{
				Containers: tt.Containers,
			}
			if err := l.Run(tt.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
