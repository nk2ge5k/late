package api

import (
	"context"

	"late"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	v1 "late/api/proto/v1"
)

var _ v1.HealthAPIServer = (*HealthService)(nil)

// HelthService is a simple service that allows to receive build information
// and check that the service is alive.
type HealthService struct{}

// Check returns build information.
func (HealthService) Check(context.Context, *emptypb.Empty) (*v1.CheckResponse, error) {
	info := late.GetBuildInfo()
	return &v1.CheckResponse{
		Version: info.Version,
		Commit:  info.Commit,
		Date:    info.Date,
	}, nil
}
