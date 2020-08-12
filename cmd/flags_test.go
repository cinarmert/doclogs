package cmd

import (
	"github.com/cinarmert/doclogs/cmd/docklogs"
	"github.com/spf13/pflag"
	"gotest.tools/v3/assert"
	"testing"
)

func Test_parseArgs(t *testing.T) {
	followFs := createFlagSetTemplate()
	err := followFs.Parse([]string{"-f"})
	assert.NilError(t, err)

	tailFs := createFlagSetTemplate()
	err = tailFs.Parse([]string{"-t5"})
	assert.NilError(t, err)

	verboseFs := createFlagSetTemplate()
	err = verboseFs.Parse([]string{"-v"})

	unknownFs := createFlagSetTemplate()
	err = unknownFs.Parse([]string{"-q"})
	assert.Check(t, err != nil, "err should not be nil for invalid args")

	invalidFs := pflag.NewFlagSet("doclogs flagset", pflag.ContinueOnError)

	type args struct {
		flags *pflag.FlagSet
		args  []string
	}
	tests := []struct {
		name    string
		args    args
		want    docklogs.LogOp
		wantErr bool
	}{
		{
			name:    "nil flagset",
			wantErr: true,
		},
		{
			name:    "empty flagset",
			args:    args{flags: createFlagSetTemplate()},
			wantErr: false,
		},
		{
			name:    "only follow in flags",
			args:    args{flags: followFs},
			wantErr: false,
		},
		{
			name:    "only tail in flags",
			args:    args{flags: tailFs},
			wantErr: false,
		},
		{
			name:    "only verbose in flags",
			args:    args{flags: verboseFs},
			wantErr: false,
		},
		{
			name:    "unknown flags in flagset",
			args:    args{flags: unknownFs},
			wantErr: false,
		},
		{
			name:    "invalid fs",
			args:    args{flags: invalidFs},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseArguments(tt.args.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArguments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func createFlagSetTemplate() *pflag.FlagSet {
	fs := pflag.NewFlagSet("doclogs flagset", pflag.ContinueOnError)
	fs.BoolP("follow", "f", false, "")
	fs.IntP("tail", "t", 100, "")
	fs.BoolP("verbose", "v", false, "")

	return fs
}
