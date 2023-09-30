// Package api contains late-service apis.
package api

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	v1 "late/api/proto/v1"
	"late/api/query"

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
			ProjectId:   req.ProjectId,
			Name:        req.Name,
			Description: req.Description,
		}
		err error
	)

	keyset, err = insertKeyset(ctx, srv.DB, keyset)
	if err != nil {
		slog.ErrorContext(ctx, "could not upsert keyset", "error", err)

		return nil, status.Errorf(codes.Internal, "could not insert keyset: %v", err)
	}

	return &v1.CreateKeysetResponse{Keyset: keyset}, nil
}

// UpdateKeyset updates keyset.
func (srv *KeysetService) UpdateKeyset(
	ctx context.Context,
	req *v1.UpdateKeysetRequest,
) (*v1.UpdateKeysetResponse, error) {
	var err error

	keyset, err := updateKeyset(ctx, srv.DB, req.Keyset)
	if err != nil {
		slog.ErrorContext(ctx, "could not upsert keyset", "error", err)

		return nil, status.Errorf(codes.Internal, "could not update keyset: %v", err)
	}

	return &v1.UpdateKeysetResponse{Keyset: keyset}, nil
}

func (srv *KeysetService) DeleteKeyset(
	ctx context.Context,
	req *v1.DeleteKeysetRequest,
) (*v1.DeleteKeysetResponse, error) {
	_, err := query.Delete_keyset_sql.Exec(ctx, srv.DB, req.Id)
	if err != nil {
		slog.ErrorContext(ctx, "could not delete keyset", "error", err)

		return nil, status.Errorf(codes.Internal, "could not delete keyset: %v", err)
	}

	return &v1.DeleteKeysetResponse{}, nil
}

func (srv *KeysetService) GetKeysets(
	ctx context.Context,
	req *v1.GetKeysetsRequest,
) (*v1.GetKeysetsResponse, error) {
	keysets, err := getKeysets(ctx, srv.DB, req.ProjectId)
	if err != nil {
		slog.ErrorContext(ctx, "could not get keysets", "error", err)

		return nil, status.Errorf(codes.Internal, "could not get keysets: %v", err)
	}

	return &v1.GetKeysetsResponse{Keysets: keysets}, nil
}

func insertKeyset(ctx context.Context, db *sql.DB, keyset *v1.Keyset) (*v1.Keyset, error) {
	row := query.Insert_keyset_sql.QueryRow(
		ctx,
		db,
		keyset.ProjectId,
		keyset.Name,
		keyset.Description,
	)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	result := new(v1.Keyset)

	err := row.Scan(&result.Id, &result.ProjectId, &result.Name, &result.Description)
	if err != nil {
		return nil, fmt.Errorf("could not scan keyset: %w", err)
	}

	return result, nil
}

func updateKeyset(ctx context.Context, db *sql.DB, keyset *v1.Keyset) (*v1.Keyset, error) {
	row := query.Update_keyset_sql.QueryRow(
		ctx,
		db,
		keyset.Id,
		keyset.Name,
		keyset.Description,
	)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	result := new(v1.Keyset)

	err := row.Scan(&result.Id, &result.ProjectId, &result.Name, &result.Description)
	if err != nil {
		return nil, fmt.Errorf("could not scan keyset: %w", err)
	}

	return result, nil
}

func getKeysets(ctx context.Context, db *sql.DB, projectID string) ([]*v1.Keyset, error) {
	rows, err := query.Get_keysets_sql.Query(ctx, db, projectID)
	if err != nil {
		return nil, fmt.Errorf("could not query rows: %w", err)
	}
	defer rows.Close()

	result := make([]*v1.Keyset, 0)

	for rows.Next() {
		keyset := new(v1.Keyset)

		err = rows.Scan(&keyset.Id, &keyset.ProjectId, &keyset.Name, &keyset.Description)
		if err != nil {
			return nil, fmt.Errorf("could not scan keyset: %w", err)
		}

		result = append(result, keyset)
	}

	return result, nil
}
