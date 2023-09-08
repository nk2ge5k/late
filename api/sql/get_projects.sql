SELECT
	id,
	name
FROM late.projects
WHERE id = ANY($1::BIGINT[]);
