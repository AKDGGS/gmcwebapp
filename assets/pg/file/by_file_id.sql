SELECT f.file_id AS id, f.description, f.mimetype AS type,
	pg_size_pretty(f.size::numeric) AS size, f.filename AS name
FROM file AS f
WHERE f.file_id = $1
