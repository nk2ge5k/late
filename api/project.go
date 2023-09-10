package api

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	v1 "late/api/proto/v1"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectService struct {
	DB *sql.DB
}

func (srv *ProjectService) CreateProject(ctx context.Context, req *v1.CreateProjectRequest) (*v1.CreateProjectResponse, error) {
	row := insert_project_sql.QueryRow(ctx, srv.DB, req.Name)
	if err := row.Err(); err != nil {
		slog.ErrorContext(ctx, "Failed to insert project",
			slog.String("error", err.Error()))
		return nil, status.Errorf(codes.Internal, "Could not insert project: %v", err)
	}

	var projectID string
	if err := row.Scan(&projectID); err != nil {
		slog.ErrorContext(ctx, "Failed to get project ID",
			slog.String("error", err.Error()))
		return nil, status.Errorf(codes.Internal, "Could not scan project's id: %v", err)
	}

	return &v1.CreateProjectResponse{
		Project: &v1.Project{
			Id:   projectID,
			Name: req.Name,
		},
	}, nil
}

func (srv *ProjectService) GetProjects(ctx context.Context, req *v1.GetProjectsRequest) (*v1.GetProjectsResponse, error) {
	projects, err := getprojects(ctx, srv.DB, req.ProjectIds...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get projects",
			slog.String("error", err.Error()))
		return nil, status.Errorf(codes.Internal, "Could not get projects: %v", err)
	}

	return &v1.GetProjectsResponse{Projects: projects}, nil
}

func getprojects(ctx context.Context, db *sql.DB, projectIDs ...string) ([]*v1.Project, error) {
	if len(projectIDs) == 0 {
		return nil, nil
	}

	ids := append([]string{}, projectIDs...)

	rows, err := get_projects_sql.Query(ctx, db, pq.StringArray(ids))
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
