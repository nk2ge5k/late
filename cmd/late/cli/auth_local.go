//go:build local
// +build local

package cli

import (
	"context"

	"late/api/auth"
)

func createAuthenticator(ctx context.Context, args *serveArgsT) (*auth.Test, error) {
	return &auth.Test{}, nil
}
