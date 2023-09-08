package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"late/api"
	adminv1 "late/api/proto/v1"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func serverCmd() *cobra.Command {

	return &cobra.Command{
		Use:   "serve",
		Short: "Start late server",
		Args:  cobra.NoArgs,
		RunE:  runServerE,
	}
}

func runServerE(cmd *cobra.Command, _ []string) error {
	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
	defer cancel()

	listener, err := net.Listen("tcp", ":0")
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
