// Package api contains late-service apis.
package api

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	v1 "late/api/proto/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ v1.KeysetAPIServer = (*KeysetService)(nil)

// KeysetService provides methods for keysets manipulations.
type KeysetService struct {
	DB *sql.DB
}

// CreateKeyset creates keyset.
func (srv *KeysetService) CreateKeyset(
	ctx context.Context,
	req *v1.CreateKeysetRequest,
) (*v1.CreateKeysetResponse, error) {
	var (
		keyset = &v1.Keyset{
			ProjectId: req.ProjectId,
			Name:      req.Name,
		}
		err error
	)

	keyset, err = upsertkeyset(ctx, srv.DB, keyset)
	if err != nil {
		slog.ErrorContext(ctx, "could not upsert keyset", "error", err)

		return nil, status.Error(codes.Internal, "could not insert keyset")
	}

	return &v1.CreateKeysetResponse{Keyset: keyset}, nil
}

// UpdateKeyset updates keyset.
func (srv *KeysetService) UpdateKeyset(
	ctx context.Context,
	req *v1.UpdateKeysetRequest,
) (*v1.UpdateKeysetResponse, error) {
	var err error

	keyset, err := upsertkeyset(ctx, srv.DB, req.Keyset)
	if err != nil {
		slog.ErrorContext(ctx, "could not upsert keyset", "error", err)

		return nil, status.Error(codes.Internal, "could not update keyset")
	}

	return &v1.UpdateKeysetResponse{Keyset: keyset}, nil
}

func (srv *KeysetService) DeleteKeyset(
	ctx context.Context,
	req *v1.DeleteKeysetRequest,
) (*v1.DeleteKeysetResponse, error) {
	_, err := delete_keyset_sql.Exec(ctx, srv.DB, req.Id)
	if err != nil {
		slog.ErrorContext(ctx, "could not delete keyset", "error", err)

		return nil, status.Error(codes.Internal, "could not delete keyset")
	}

	return &v1.DeleteKeysetResponse{}, nil
}

func (srv *KeysetService) GetKeysets(
	ctx context.Context,
	req *v1.GetKeysetsRequest,
) (*v1.GetKeysetsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

func upsertkeyset(ctx context.Context, db *sql.DB, keyset *v1.Keyset) (*v1.Keyset, error) {
	row := upsert_keyset_sql.QueryRow(
		ctx,
		db,
		keyset.ProjectId,
		keyset.Name,
	)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	result := new(v1.Keyset)

	err := row.Scan(&result.Id, &result.ProjectId, &result.Name)
	if err != nil {
		return nil, fmt.Errorf("could not scan keyset: %w", err)
	}

	return result, nil
}
