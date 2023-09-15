CREATE SEQUENCE keysets_id_seq;

CREATE TABLE late.keysets (
    id          TEXT NOT NULL DEFAULT nextval('keysets_id_seq')::TEXT PRIMARY KEY,
    project_id  TEXT NOT NULL REFERENCES projects(id),
    name        TEXT NOT NULL,
	description TEXT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
