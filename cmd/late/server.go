package main

import "github.com/spf13/cobra"

func serverCmd() *cobra.Command {

	return &cobra.Command{
		Use:   "serve",
		Short: "Start late server",
		Args:  cobra.NoArgs,
		RunE:  runServerE,
	}
}

func runServerE(cmd *cobra.Command, _ []string) error {
	return nil
}
