CREATE TABLE late.keys (
    keyset_id    TEXT NOT NULL REFERENCES late.keysets(id),
    key          TEXT NOT NULL,
    description  TEXT NULL,
    translations JSONB NOT NULL DEFAULT '[]',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),


    UNIQUE(keyset_id, key)
);
