INSERT INTO file (
	file_id, filename, description, size, mimetype, file_hash, file_type
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
)
