package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffyaml"
)

const rootCommandName = "late"

// Stderr and Stdout aliases for test purposes.
var Stderr io.Writer = os.Stderr
var Stdout io.Writer = os.Stdout

// Run runs the CLI. The args do not include the binary name.
func Run(args []string) error {
	rootCmd := &ff.Command{
		Name:      rootCommandName,
		ShortHelp: fmt.Sprintf("%s [flags] <subcommand> [command flags]", rootCommandName),
		Flags:     newFlagSet(rootCommandName),
		Exec:      func(context.Context, []string) error { return flag.ErrHelp },
		Subcommands: []*ff.Command{
			serveCmd,
			versionCmd,
		},
	}

	parser := ffyaml.Parser{Delimiter: "."}

	if err := rootCmd.Parse(args,
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(parser.Parse)); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := rootCmd.Run(ctx); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	return nil
}

func newFlagSet(name string) *ff.FlagSet {
	fs := ff.NewFlagSet(name)
	fs.String('c', "config", "", "config file")

	return fs
}
