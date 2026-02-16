UPDATE file SET
	filename = $1,
	description = $2,
	size = $3,
	mimetype = $4,
	file_hash = $5
WHERE file_id = $6
