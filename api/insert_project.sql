INSERT INTO projects (name)
VALUES ($1)
RETURNING id;
