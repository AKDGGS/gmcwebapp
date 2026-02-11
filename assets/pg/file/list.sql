SELECT
	f.file_id AS id,
	f.filename AS name,
	f.description,
	f.mimetype,
	f.size::numeric,
    f.file_hash
FROM file AS f