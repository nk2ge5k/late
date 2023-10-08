UPDATE late.keys SET
	description = $3,
	translations = $4
WHERE keyset_id = $1 AND key = $2;
