CREATE SEQUENCE projects_id_seq;

CREATE TABLE projects (
  id TEXT NOT NULL DEFAULT nextval('projects_id_seq')::TEXT PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
