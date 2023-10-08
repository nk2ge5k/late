package api

import (
	"context"
	"database/sql"
	"encoding/json"

	v1 "late/api/proto/v1"
	"late/api/query"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ v1.KeysAPIServer = (*KeysService)(nil)

// KeysService implements gRPC Keys API server.
type KeysService struct {
	DB *sql.DB
}

func (srv *KeysService) CreateKey(
	ctx context.Context,
	req *v1.CreateKeyRequest,
) (*v1.CreateKeyResponse, error) {
	// NOTE(vitaminniy): we insert raw json into translations cause we don't
	// have enough mental power to fight pg arrays.
	rawtranslations, err := json.Marshal(req.Translations)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not marshal translations: %v", err)
	}

	_, err = query.Insert_key_sql.Exec(
		ctx,
		srv.DB,
		req.KeysetId,
		req.Key,
		req.Description,
		rawtranslations,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not exec insert: %v", err)
	}

	return &v1.CreateKeyResponse{}, nil
}

func (srv *KeysService) UpdateKey(
	ctx context.Context,
	req *v1.UpdateKeyRequest,
) (*v1.UpdateKeyResponse, error) {
	// NOTE(vitaminniy): we insert raw json into translations cause we don't
	// have enough mental power to fight pg arrays.
	rawtranslations, err := json.Marshal(req.Translations)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not marshal translations: %v", err)
	}

	_, err = query.Update_key_sql.Exec(
		ctx,
		srv.DB,
		req.KeysetId,
		req.Key,
		req.Description,
		rawtranslations,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not exec query: %v", err)
	}

	return &v1.UpdateKeyResponse{}, nil
}

func (srv *KeysService) DeleteKey(
	ctx context.Context,
	req *v1.DeleteKeyRequest,
) (*v1.DeleteKeyResponse, error) {
	_, err := query.Delete_key_sql.Exec(ctx, srv.DB, req.KeysetId, req.Key)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not exec query: %v", err)
	}

	return &v1.DeleteKeyResponse{}, nil
}

func (srv *KeysService) GetKeys(
	ctx context.Context,
	req *v1.GetKeysRequest,
) (*v1.GetKeysResponse, error) {
	rows, err := query.Get_keys_sql.Query(ctx, srv.DB, req.KeysetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not query rows: %v", err)
	}
	defer rows.Close()

	var keys []*v1.Key

	for rows.Next() {
		key := new(v1.Key)
		translations := make([]byte, 0)

		err := rows.Scan(
			&key.KeysetId,
			&key.Key,
			&key.Description,
			&translations,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not scan key: %v", err)
		}

		if err := json.Unmarshal(translations, &key.Translations); err != nil {
			return nil, status.Errorf(
				codes.Internal, "could not unmarshal translations: %v", err)
		}

		keys = append(keys, key)
	}

	return &v1.GetKeysResponse{Keys: keys}, nil
}
