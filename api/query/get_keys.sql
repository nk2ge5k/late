SELECT
	keyset_id,
	key,
	description,
	translations
FROM late.keys
WHERE keyset_id = $1;
