SELECT f.file_id AS id,
	f.filename AS name,
	f.description,
	f.mimetype AS type,
	f.size::numeric AS size,
	f.content_md5 AS MD5
FROM file AS f
WHERE f.file_id = $1
