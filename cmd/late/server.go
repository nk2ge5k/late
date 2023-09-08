package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"

	"late/api"
	adminv1 "late/api/proto/v1"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var configfile string

func serverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start late server",
		Args:  cobra.NoArgs,
		RunE:  runServerE,
	}

	cmd.Flags().StringVarP(&configfile, "config", "c", "", "path to config file; must be set")

	return cmd
}

func runServerE(cmd *cobra.Command, args []string) error {
	if err := cmd.Flags().Parse(args); err != nil {
		return fmt.Errorf("could not parse args: %w", err)
	}

	if configfile == "" {
		return errors.New("'config' flag must be set")
	}

	cfg, err := LoadConfig(configfile)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
	defer cancel()

	db, err := sql.Open("postgres", cfg.Postgres.URL)
	if err != nil {
		return fmt.Errorf("could not open postgres: %w", err)
	}
	defer db.Close()

	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("could not establish db conn: %w", err)
	}

	listener, err := net.Listen("tcp", cfg.GRPC.Addr)
	if err != nil {
		return fmt.Errorf("could not start listener: %w", err)
	}
	defer listener.Close()

	var service api.ProjectService

	srv := grpc.NewServer()
	srv.RegisterService(&adminv1.ProjectAPI_ServiceDesc, &service)

	errs := make(chan error, 1)

	go func() {
		errs <- srv.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		srv.GracefulStop()
	case err := <-errs:
		return fmt.Errorf("could not serve gRPC: %w", err)
	}

	return nil
}
