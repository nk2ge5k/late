package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"late/api"
	v1 "late/api/proto/v1"
	"late/metrics"
	"late/shutdown"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	srv := grpc.NewServer()
	{
		db, err := sql.Open("postgres", cfg.Postgres.URI)
		if err != nil {
			return fmt.Errorf("could not open postgres: %w", err)
		}
		defer db.Close()

		if err = db.PingContext(ctx); err != nil {
			return fmt.Errorf("could not establish db conn: %w", err)
		}

		v1.RegisterProjectAPIServer(srv, &api.ProjectService{DB: db})
		v1.RegisterHealthAPIServer(srv, api.HealthService{})
	}

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		shutdown.Handle(srv.GracefulStop)

		ln, err := net.Listen("tcp", cfg.GRPC.Listen)
		if err != nil {
			return fmt.Errorf("could not start listener: %w", err)
		}
		defer ln.Close()

		slog.Info("Starting GRPC server", slog.String("addr", cfg.GRPC.Listen))
		if err := srv.Serve(ln); err != nil {
			return fmt.Errorf("grpc serve: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		metricsHandler := metrics.Handler()

		mux := runtime.NewServeMux()
		//nolint:errcheck
		mux.HandlePath(http.MethodGet, "/metrics", runtime.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request, _ map[string]string) {
				metricsHandler.ServeHTTP(rw, req)
			}))
		rerr := reigisterGRPCGateways(gctx, mux, cfg.GRPC.Listen,
			v1.RegisterHealthAPIHandlerFromEndpoint)
		if rerr != nil {
			return rerr
		}

		srv := http.Server{
			Addr:              cfg.HTTP.Listen,
			Handler:           mux,
			ReadHeaderTimeout: time.Millisecond,
		}

		slog.Info("Starting HTTP server", slog.String("addr", cfg.HTTP.Listen))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http listen and serve: %w", err)
		}
		return nil
	})

	return g.Wait()
}

type RegisterFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

func reigisterGRPCGateways(
	ctx context.Context, mux *runtime.ServeMux, endpoint string, fns ...RegisterFunc) error {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	for _, fn := range fns {
		if err := fn(ctx, mux, endpoint, opts); err != nil {
			return fmt.Errorf("register grpc gateway: %w", err)
		}
	}
	return nil
}
