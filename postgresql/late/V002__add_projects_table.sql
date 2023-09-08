CREATE SEQUENCE late.projects_id_seq;

CREATE TABLE late.projects AS (
	id TEXT NOT NULL DEFAULT CAST(nextval('late.projects_id_set') AS TEXT),
	name TEXT NOT NULL UNIQUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER SEQUENCE late.projects_id_seq OWNED BY late.projects;
