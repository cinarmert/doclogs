package cmd

import (
	"github.com/cinarmert/doclogs/cmd/docklogs"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

// parseArguments parses the provided args and flags and return the
// corresponding operation.
func parseArguments(flags *pflag.FlagSet, args ...string) (Op, error) {
	if flags == nil {
		return nil, errors.New("flagset cannot be nil")
	}

	follow, err := flags.GetBool("follow")
	if err != nil {
		return nil, errors.Wrap(err, "could not parse follow flag")
	}

	return &docklogs.LogOp{
		Follow:     follow,
		Containers: args,
	}, nil
}
