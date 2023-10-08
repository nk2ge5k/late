DELETE FROM late.keys
WHERE keyset_id = $1 AND key = $2;
