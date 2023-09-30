package auth

import (
	"context"

	fireauth "firebase.google.com/go/v4/auth"
)

type Firebase struct {
	cli *fireauth.Client
}

// New returns new instance of firebase authenticator.
func New(cli *fireauth.Client) *Firebase {
	return &Firebase{cli: cli}
}

// VerifyToken verifies the provided ID token, and additionally checks that the
// token has not been revoked or disabled.
func (f *Firebase) VerifyToken(ctx context.Context, tok string) (*Token, error) {
	t, err := f.cli.VerifyIDTokenAndCheckRevoked(ctx, tok)
	if err != nil {
		return nil, err
	}

	return &Token{
		AuthTime: t.AuthTime,
		Issuer:   t.Issuer,
		Audience: t.Audience,
		Expires:  t.Expires,
		IssuedAt: t.IssuedAt,
		Subject:  t.Subject,
		UID:      t.UID,
	}, nil
}
