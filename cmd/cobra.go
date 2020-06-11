package cmd

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

type Op interface {
	Run(context.Context) error
}

var rootCmd = &cobra.Command{
	Use:   "doclogs [OPTIONS] [CONTAINERS...]",
	Short: "Multiple Docker Container Log Visualizer",
	Long:  `Doclogs provides a user interface for the logs from multiple docker containers.`,
	Run:   run,
}

func init() {
	rootCmd.PersistentFlags().BoolP("follow", "f", false, "follow the stream of logs")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "print debug logs")
	log.SetLevel(log.WarnLevel)
}

func run(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	if len(args) == 0 {
		fmt.Println("No docker container provided\n ")
		cmd.Help()
		os.Exit(1)
	}

	op, err := parseArguments(cmd.Flags(), args...)
	if err != nil {
		log.Errorf("could not parse given flags: %v", err)
		os.Exit(1)
	}

	if err := op.Run(ctx); err != nil {
		log.Errorf("could not get container logs, use -v flag to enable verbose logs")

		if v, _ := cmd.Flags().GetBool("verbose"); v {
			log.Errorf("[DEBUG]: %+v\n", err)
		}
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error running the init command")
	}
}
