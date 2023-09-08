package api

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"io"

	v1 "late/api/proto/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	//go:embed sql/insert_project.sql
	insertProjectQuery string
	getProjectQuery    string
)

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

func (srv *ProjectService) GetProject(stream v1.ProjectAPI_GetProjectServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		project, err := getproject(stream.Context(), srv.DB, req.Id)
		if err != nil {
			return status.Errorf(codes.Internal, "could not get project")
		}

		if project == nil {
			continue
		}

		if err = stream.Send(&v1.GetProjectResponse{Project: project}); err != nil {
			return err
		}
	}
}

func (srv *ProjectService) GetProjects(context.Context, *v1.GetProjectsRequest) (*v1.GetProjectsResponse, error) {
	return &v1.GetProjectsResponse{}, nil
}

func getproject(ctx context.Context, db *sql.DB, projectID string) (*v1.Project, error) {
	row := db.QueryRowContext(ctx, getProjectQuery, projectID)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("could not query row: %w", err)
	}

	proj := new(v1.Project)
	if err := row.Scan(&proj.Id, &proj.Name); err != nil {
		return nil, fmt.Errorf("could not scan row: %w", err)
	}

	return proj, nil
}
