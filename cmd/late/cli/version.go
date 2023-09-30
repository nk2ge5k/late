package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"late"

	"github.com/peterbourgon/ff/v4"
)

const versionCommandName = "version"

var versionCmd = &ff.Command{
	Name:      versionCommandName,
	ShortHelp: fmt.Sprintf("%s [flags]", versionCommandName),
	Flags:     versionFlagSet,
	Exec: func(context.Context, []string) error {
		info := late.GetBuildInfo()

		if versionArgsGlobal.json {
			enc := json.NewEncoder(Stdout)
			enc.SetIndent("", " ")

			if err := enc.Encode(info); err != nil {
				return fmt.Errorf("json encode: %w", err)
			}

			return nil
		}

		tw := tabwriter.NewWriter(Stdout, 0, 2, 2, ' ', 0)

		fmt.Fprintf(tw, "Version\t%s\n", info.Version)
		fmt.Fprintf(tw, "Commit\t%s\n", info.Commit)
		fmt.Fprintf(tw, "Date\t%s\n", info.Date)

		tw.Flush()

		return nil
	},
}

// versionArgsT contains command-line arguments values
type versionArgsT struct {
	json bool
}

var (
	versionArgsGlobal versionArgsT
	versionFlagSet    = newVersionFlagSet(&versionArgsGlobal)
)

// newVersionFlagSet returns new flagset for serve command.
func newVersionFlagSet(versionArgs *versionArgsT) *ff.FlagSet {
	fs := newFlagSet(versionCommandName)
	fs.BoolVarDefault(&versionArgs.json, 'j', "json", false, "Output as JSON")
	return fs
}
