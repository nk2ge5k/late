SELECT
  id,
  name
FROM
  projects
WHERE
  id = ANY ($1);
