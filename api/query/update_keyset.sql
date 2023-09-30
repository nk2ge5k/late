UPDATE late.keysets SET
	name = $2,
	description = $3,
	updated_at = now()
WHERE id = $1;
