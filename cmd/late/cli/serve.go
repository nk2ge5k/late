package cli

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"late"

	stdruntime "runtime"

	"late/api"
	apiauth "late/api/auth"
	v1 "late/api/proto/v1"
	"late/metrics"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	_ "github.com/lib/pq"
	"github.com/peterbourgon/ff/v4"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const serveCommandName = "serve"

var serveCmd = &ff.Command{
	Name:      serveCommandName,
	ShortHelp: fmt.Sprintf("%s [flags]", serveCommandName),
	Flags:     serveFlagSet,
	Exec: func(ctx context.Context, args []string) error {
		return runServe(ctx, args, serveArgsGlobal)
	},
}

// serveArgsT contains command-line arguments values
type serveArgsT struct {
	Postgres struct {
		URI string `yaml:"uri"`
	} `yaml:"postgres"`
	GRPC struct {
		Listen string `yaml:"listen"`
	} `yaml:"grpc"`
	HTTP struct {
		Listen string `yaml:"listen"`
	} `yaml:"http"`
	Firebase struct {
		ProjectID       string `yaml:"project-id"`
		CredentialsFile string `yaml:"credentials-file"`
		ServiceAccout   string `yaml:"service-accout-file"`
	} `yaml:"firebase"`
}

var (
	serveArgsGlobal serveArgsT
	serveFlagSet    = newserveFlagSet(&serveArgsGlobal)
)

// newserveFlagSet returns new flagset for serve command.
func newserveFlagSet(serveArgs *serveArgsT) *ff.FlagSet {
	fs := newFlagSet(serveCommandName)
	fs.StringVar(&serveArgs.Postgres.URI,
		0, "postgres.uri", "", "Postgres connection string")
	fs.StringVar(&serveArgs.GRPC.Listen,
		0, "grpc.listen", "localhost:18443",
		"address on which to listen for remote connections to the GPRC server")
	fs.StringVar(&serveArgs.HTTP.Listen,
		0, "http.listen", "localhost:18080",
		"address on which to listen for remote connections to the HTTP server")
	fs.StringVar(&serveArgs.Firebase.ProjectID,
		0, "firebase.project-id", "", "Firbase project id")
	fs.StringVar(&serveArgs.Firebase.CredentialsFile,
		0, "firebase.credentials-file", "",
		"Path to file containing firebase application credentials")
	return fs
}

var (
	logInterceptor = logging.LoggerFunc(func(
		ctx context.Context, level logging.Level, msg string, fields ...any) {
		slog.Log(ctx, slog.Level(level), msg, fields...)
	})
	panicInterceptor = recovery.RecoveryHandlerFunc(func(p any) (err error) {
		stack := make([]byte, 64<<10)
		stack = stack[:stdruntime.Stack(stack, false)]

		slog.Error("GRPC handler panic",
			slog.Any("data", p), slog.String("stack", string(stack)))

		return status.Error(codes.Internal, "internal server error")
	})
)

//nolint:gocritic
func runServe(ctx context.Context, _ []string, args serveArgsT) error {
	authenticator, err := createAuthenticator(ctx, &args)
	if err != nil {
		return err
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(logInterceptor),
			auth.UnaryServerInterceptor(apiauth.GRPCInterceptor(authenticator)),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(panicInterceptor)),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(logInterceptor),
			auth.StreamServerInterceptor(apiauth.GRPCInterceptor(authenticator)),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(panicInterceptor)),
		),
	)
	{
		db, err := sql.Open("postgres", args.Postgres.URI)
		if err != nil {
			return fmt.Errorf("could not open postgres: %w", err)
		}
		defer db.Close()

		if err = db.PingContext(ctx); err != nil {
			return fmt.Errorf("could not establish db conn: %w", err)
		}

		v1.RegisterHealthAPIServer(srv, api.HealthService{})
		v1.RegisterProjectAPIServer(srv, &api.ProjectService{DB: db})
		v1.RegisterKeysetAPIServer(srv, &api.KeysetService{DB: db})
	}

	g, gctx := errgroup.WithContext(ctx)
	g.Go(grpcServerRunner(gctx, srv, &args))
	g.Go(httpServerRunner(gctx, srv, &args))

	return g.Wait()
}

// grpcServerRunner return runner for the GRPC server
func grpcServerRunner(ctx context.Context, srv *grpc.Server, cfg *serveArgsT) func() error {
	return func() error {
		ln, err := net.Listen("tcp", cfg.GRPC.Listen)
		if err != nil {
			return fmt.Errorf("could not start listener: %w", err)
		}
		defer ln.Close()

		slog.Info("Starting GRPC server", slog.String("addr", cfg.GRPC.Listen))

		go func() {
			<-ctx.Done()
			slog.Info("Stopping GRPC server")
			srv.GracefulStop()
		}()

		if err := srv.Serve(ln); err != nil {
			return fmt.Errorf("grpc serve: %w", err)
		}
		return nil
	}
}

// httpServerRunner returns runner for the HTTP server
func httpServerRunner(ctx context.Context, srv *grpc.Server, cfg *serveArgsT) func() error {
	return func() error {
		var (
			metricsHandler = metrics.Handler()
			grpcHandler    = grpcweb.WrapServer(srv)
		)

		srv := http.Server{
			Addr: cfg.HTTP.Listen,
			Handler: http.HandlerFunc(
				func(rw http.ResponseWriter, req *http.Request) {
					switch {
					case grpcHandler.IsGrpcWebRequest(req):
						grpcHandler.ServeHTTP(rw, req)
					case req.Method == "GET" && strings.HasPrefix(req.URL.Path, "/metrics"):
						metricsHandler.ServeHTTP(rw, req)
					case req.Method == "GET" && strings.HasPrefix(req.URL.Path, "/healthz"):
						b, err := json.Marshal(late.GetBuildInfo())
						if err != nil {
							http.Error(rw, err.Error(), http.StatusInternalServerError)
							return
						}
						rw.WriteHeader(http.StatusOK)
						rw.Write(b) //nolint:errcheck
					default:
						http.Error(rw, fmt.Sprint(req.Method, req.URL.Path, "not found"),
							http.StatusNotFound)
					}
				},
			),
			ReadHeaderTimeout: time.Millisecond,
		}

		go func() {
			<-ctx.Done()
			slog.Info("Stopping HTTP server")
			if err := srv.Shutdown(context.Background()); err != nil {
				slog.Error("Failed to shutdown server",
					slog.String("error", err.Error()))
			}
		}()

		slog.Info("Starting HTTP server", slog.String("addr", cfg.HTTP.Listen))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http listen and serve: %w", err)
		}
		return nil
	}
}
