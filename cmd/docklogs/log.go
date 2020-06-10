package docklogs

import (
	"context"
	"fmt"
	"github.com/cinarmert/doclogs/cmd/docklogs/container"
	"github.com/cinarmert/doclogs/cmd/docklogs/ui"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type LogOp struct {
	Follow     bool
	FileOut    bool
	Containers []string
}

// Run is the entrypoint for fetching the logs
func (l *LogOp) Run(ctx context.Context) error {
	if ctx == nil {
		return errors.New("context parameter cannot be nil")
	}

	dockerClient, err := GetDockerClient(ctx)
	if err != nil {
		return errors.Wrap(err, "could not initiate docker client")
	}

	sessions := l.createContainerSessions(ctx, dockerClient)

	if len(sessions) == 0 {
		return errors.New("could not create any container session (probably because none are running), exiting")
	}

	lm, err := ui.NewManagerForSessions(sessions)

	if err != nil {
		return errors.Wrap(err, "could not create layout manager")
	}
	lm.Run()

	l.report(sessions)
	return nil
}

// report is executed when the process is about to end to
// inform about the condition of the given containers.
func (l *LogOp) report(sessions []*container.Session) {
	for _, s := range sessions {
		fmt.Printf("container %s status: %s\n", s.Name, s.Status.String())
	}
}

// createContainerSessions creates sessions with given container names.
// It silently discards the dead containers with given name.
func (l *LogOp) createContainerSessions(ctx context.Context, client DockerCli) []*container.Session {
	var containers []*container.Session
	for _, name := range l.Containers {
		if client.IsContainerLive(name) {

			tmp := container.NewContainerSession(client, ctx, name, l.Follow)
			containers = append(containers, tmp)
		} else {
			log.Warnf("container \"%s\" does not appear to be live, skipping", name)
		}
	}
	return containers
}
