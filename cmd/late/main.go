package main

import (
	"fmt"
	"late"
	"late/shutdown"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "late"}
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(serverCmd())

	if err := rootCmd.Execute(); err != nil {
		handleErr(err.Error())
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the server version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "Late %v\n", late.GetBuildInfo())
		},
	}
}

func handleErr(err string) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	shutdown.Now(1)
}
