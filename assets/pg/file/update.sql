UPDATE file SET
	filename = COALESCE($2, filename),
	description = COALESCE($3, description),
	size = COALESCE($4, size),
	mimetype = COALESCE($5, mimetype),
	file_hash = COALESCE($6, file_hash)
WHERE file_id = $1
