//go:build !local
// +build !local

package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	apiauth "late/api/auth"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func createAuthenticator(ctx context.Context, args *serveArgsT) (*apiauth.Firebase, error) {
	if args.Firebase.ProjectID == "" {
		return nil, errors.New("missing firebase.project-id")
	}
	if args.Firebase.CredentialsFile == "" || args.Firebase.ServiceAccout == "" {
		return nil, errors.New(
			"either firebase.credentials-file or firebase.service-accout must be provided")
	}

	opts := []option.ClientOption{
		option.WithCredentialsFile(args.Firebase.CredentialsFile),
	}
	if host, ok := os.LookupEnv("FIREBASE_AUTH_EMULATOR_HOST"); ok {
		slog.Warn("Using firebase emulator host", slog.String("host", host))
		opts = append(opts, option.WithEndpoint(host))
	}

	app, err := firebase.NewApp(
		ctx, &firebase.Config{ProjectID: args.Firebase.ProjectID}, opts...)
	if err != nil {
		return nil, fmt.Errorf("firebase new app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("firebase auth: %w", err)
	}

	return apiauth.New(client), nil
}
