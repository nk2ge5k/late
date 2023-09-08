package api

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"io"

	v1 "late/api/proto/v1"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	//go:embed sql/insert_project.sql
	insertProjectQuery string
	//go:embed sql/get_projects.sql
	getProjectsQuery string
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

		projects, err := getprojects(stream.Context(), srv.DB, req.ProjectId)
		if err != nil {
			return status.Errorf(codes.Internal, "could not get project")
		}

		if len(projects) == 0 {
			continue
		}

		response := &v1.GetProjectResponse{Project: projects[0]}

		if err = stream.Send(response); err != nil {
			return err
		}
	}
}

func (srv *ProjectService) GetProjects(ctx context.Context, req *v1.GetProjectsRequest) (*v1.GetProjectsResponse, error) {
	projects, err := getprojects(ctx, srv.DB, req.ProjectIds...)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not get projects")
	}

	return &v1.GetProjectsResponse{Projects: projects}, nil
}

func getprojects(ctx context.Context, db *sql.DB, projectIDs ...string) ([]*v1.Project, error) {
	if len(projectIDs) == 0 {
		return nil, nil
	}

	ids := append([]string{}, projectIDs...)

	rows, err := db.QueryContext(ctx, getProjectsQuery, pq.StringArray(ids))
	if err != nil {
		return nil, fmt.Errorf("could not query rows: %w", err)
	}
	defer rows.Close()

	projects := make([]*v1.Project, 0, len(projectIDs))

	for rows.Next() {
		project := new(v1.Project)

		if err := rows.Scan(&project.Id, &project.Name); err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}

		projects = append(projects, project)
	}

	return projects, nil
}
