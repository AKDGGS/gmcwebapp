SELECT 
	f.file_id AS id,
	f.filename AS name,
	f.description,
	f.mimetype,
	f.size::numeric,
	f.content_md5 AS MD5
FROM file AS f
WHERE f.file_id = $1
