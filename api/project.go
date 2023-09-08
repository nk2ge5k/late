package api

import (
	"context"
	"database/sql"
	_ "embed"

	v1 "late/api/proto/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed sql/insert_project.sql
var insertProjectQuery string

type ProjectService struct {
	DB *sql.DB
}

func (srv *ProjectService) CreateProject(ctx context.Context, req *v1.CreateProjectRequest) (*v1.CreateProjectResponse, error) {
	row := srv.DB.QueryRowContext(ctx, insertProjectQuery, req.Name)
	if err := row.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert project")
	}

	var projectID string
	if err := row.Scan(&projectID); err != nil {
		return nil, status.Errorf(codes.Internal, "could not scan project's id")
	}

	return &v1.CreateProjectResponse{
		Project: &v1.Project{
			Id:   projectID,
			Name: req.Name,
		},
	}, nil
}

func (srv *ProjectService) GetProject(v1.ProjectAPI_GetProjectServer) error {
	return status.Errorf(codes.Unimplemented, "method GetProject not implemented")
}

func (srv *ProjectService) GetProjects(context.Context, *v1.GetProjectsRequest) (*v1.GetProjectsResponse, error) {
	return &v1.GetProjectsResponse{}, nil
}
