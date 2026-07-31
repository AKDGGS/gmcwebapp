SELECT
	f.file_id AS id,
	f.filename AS name,
	f.file_hash,
	f.description,
	f.mimetype,
	f.size::numeric
FROM file AS f
WHERE f.file_hash = $1
