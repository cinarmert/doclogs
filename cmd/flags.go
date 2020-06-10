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

	out, err := flags.GetBool("output")
	if err != nil {
		return nil, errors.Wrap(err, "could not parse output flag")
	}

	if follow && out {
		return nil, errors.New("output and follow cannot be given at the same time")
	}

	return &docklogs.LogOp{
		Follow:     follow,
		FileOut:    out,
		Containers: args,
	}, nil
}
