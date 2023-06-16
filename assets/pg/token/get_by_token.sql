SELECT token_id as id, description
	FROM token
	WHERE token = $1
