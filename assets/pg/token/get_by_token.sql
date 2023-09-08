SELECT token_id AS id,
	description
FROM token
WHERE token = $1
