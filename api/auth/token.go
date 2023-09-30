package auth

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Token represents a decoded ID token.
type Token struct {
	AuthTime int64  `json:"auth_time"`
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
	Expires  int64  `json:"exp"`
	IssuedAt int64  `json:"iat"`
	Subject  string `json:"sub,omitempty"`
	UID      string `json:"uid,omitempty"`
}

var tokenKey struct{}

// WithToken returns new context with token value.
func WithToken(ctx context.Context, tok *Token) context.Context {
	return context.WithValue(ctx, tokenKey, tok)
}

// ExtractToken attempts to extract token form given context.
func ExtractToken(ctx context.Context) (tok *Token, ok bool) {
	tok, ok = ctx.Value(tokenKey).(*Token)
	return tok, ok
}

// GRPCInterceptor returns AuthFunc for GRPC server interceptor.
func GRPCInterceptor(verifier interface {
	VerifyToken(ctx context.Context, tok string) (*Token, error)
}) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		idToken, err := auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		tok, verr := verifier.VerifyToken(ctx, idToken)
		if verr != nil {
			return nil, status.Errorf(codes.Unauthenticated,
				"Token verifiction: %v", verr)
		}

		return WithToken(ctx, tok), nil
	}
}
