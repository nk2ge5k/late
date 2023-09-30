//go:build local
// +build local

package auth

import (
	"context"
	"time"
)

type Test struct{}

// VerifyToken verifies the provided ID token, and additionally checks that the
// token has not been revoked or disabled.
func (t *Test) VerifyToken(ctx context.Context, tok string) (*Token, error) {
	// TODO (nk2ge5k): It might be useful to add some capabilities to this test
	//	implementation that would allow adding and removing users for testing
	//	purposes.
	now := time.Now()
	return &Token{
		AuthTime: now.Unix(),
		Issuer:   "testsuite",
		Audience: "late",
		Expires:  now.Add(time.Hour).Unix(),
		IssuedAt: now.Unix(),
		Subject:  "late",
		UID:      "testsuite",
	}, nil
}
