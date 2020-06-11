package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"testing"
)

func Test_runCrashOnEmptyArgs(t *testing.T) {
	c := &cobra.Command{}
	if os.Getenv("SUBPROCESS_CRASHER") == "1" {
		run(c, nil)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=Test_runCrashOnEmptyArgs")
	cmd.Env = append(os.Environ(), "SUBPROCESS_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func Test_runCrashOnInvalidCobraCommand(t *testing.T) {
	c := &cobra.Command{}
	if os.Getenv("SUBPROCESS_CRASHER") == "1" {
		run(c, []string{"aa"})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=Test_runCrashOnInvalidCobraCommand")
	cmd.Env = append(os.Environ(), "SUBPROCESS_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func Test_runCrashOnDockerNotRunning(t *testing.T) {
	c := &cobra.Command{}
	c.PersistentFlags().BoolP("follow", "f", false, "follow the stream of logs")
	c.PersistentFlags().BoolP("verbose", "v", false, "print debug logs")
	if os.Getenv("SUBPROCESS_CRASHER") == "1" {
		run(c, []string{"aa"})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=Test_runCrashOnInvalidCobraCommand")
	cmd.Env = append(os.Environ(), "SUBPROCESS_CRASHER=1", "_FORCE_DOCKER_RUNNING=False")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
