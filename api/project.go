package api

import (
	"context"

	v1 "late/api/proto/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectService struct {
}

func (srv *ProjectService) CreateProject(context.Context, *v1.CreateProjectRequest) (*v1.CreateProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProject not implemented")
}

func (srv *ProjectService) GetProject(v1.ProjectAPI_GetProjectServer) error {
	return status.Errorf(codes.Unimplemented, "method GetProject not implemented")
}

func (srv *ProjectService) GetProjects(context.Context, *v1.GetProjectsRequest) (*v1.GetProjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProjects not implemented")
}
