package ui

import (
	"context"
	"github.com/cinarmert/doclogs/cmd/docklogs/container"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/rivo/tview"
	"gotest.tools/assert"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"
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

func TestLayoutManager_Run(t *testing.T) {
	f := &fakeClient{
		logFunc: func(s string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
			return ioutil.NopCloser(strings.NewReader("test log")), nil
		},
	}

	s := container.NewContainerSession(f, context.TODO(), "test", false)

	tests := []struct {
		name     string
		sessions []*container.Session
	}{
		{
			name:     "no session given",
			sessions: []*container.Session{},
		},
		{
			name:     "single sessions given",
			sessions: []*container.Session{s},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			lm := &LayoutManager{
				App:      tview.NewApplication(),
				Grid:     tview.NewGrid(),
				Sessions: tt.sessions,
			}
			go lm.Run()
			time.Sleep(time.Millisecond * 300)
		})
	}
}

func TestLayoutManager_SetApp(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful call",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := &LayoutManager{}
			got := lm.SetApp()
			assert.Check(t, got.App != nil)
		})
	}
}

func TestLayoutManager_SetGrid(t *testing.T) {
	tests := []struct {
		name     string
		sessions []*container.Session
		wantLen  int
	}{
		{
			name:     "no session given",
			sessions: []*container.Session{},
			wantLen:  0,
		},
		{
			name:     "multiple sessions given",
			sessions: []*container.Session{{Name: "test1"}, {Name: "test2"}, {Name: "test3"}, {Name: "test4"}},
			wantLen:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := &LayoutManager{
				App:      tview.NewApplication(),
				Grid:     tview.NewGrid(),
				Sessions: tt.sessions,
			}
			got := lm.SetGrid()
			assert.Check(t, got != nil)
			assert.Check(t, lm.Grid != nil)
		})
	}
}

func TestLayoutManager_SetSessions(t *testing.T) {
	tests := []struct {
		name     string
		sessions []*container.Session
		wantLen  int
	}{
		{
			name:     "no session given",
			sessions: []*container.Session{},
			wantLen:  0,
		},
		{
			name:     "single session given",
			sessions: []*container.Session{{Name: "test1"}},
			wantLen:  1,
		},
		{
			name:     "multiple sessions given",
			sessions: []*container.Session{{Name: "test1"}, {Name: "test2"}, {Name: "test3"}, {Name: "test4"}},
			wantLen:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := &LayoutManager{
				App:  tview.NewApplication(),
				Grid: tview.NewGrid(),
			}
			got := lm.SetSessions(tt.sessions)
			assert.Equal(t, len(got.Sessions), tt.wantLen)
		})
	}
}

func TestLayoutManager_createTextView(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantTitle string
	}{
		{
			name:      "empty title",
			title:     "",
			wantTitle: "  ",
		},
		{
			name:      "regular title",
			title:     "test title",
			wantTitle: " test title ",
		},
		{
			name:      "special chars in title",
			title:     "!/#test\"",
			wantTitle: " !/#test\" ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := &LayoutManager{
				App:  tview.NewApplication(),
				Grid: tview.NewGrid(),
			}
			got := lm.createTextView(tt.title)
			assert.Equal(t, got.GetTitle(), tt.wantTitle)
		})
	}
}

func TestNewManagerForSessions(t *testing.T) {
	type args struct {
		sessions []*container.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful creation",
			args: args{
				sessions: []*container.Session{{Name: "test"}},
			},
			wantErr: false,
		},
		{
			name: "successful creation",
			args: args{
				sessions: []*container.Session{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewManagerForSessions(tt.args.sessions)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManagerForSessions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
