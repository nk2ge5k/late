CREATE SEQUENCE projects_id_seq;

CREATE TABLE projects (
  id          TEXT NOT NULL DEFAULT nextval('projects_id_seq')::TEXT,
	name        TEXT NOT NULL,
	created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
