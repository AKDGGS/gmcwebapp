SELECT
	f.file_id AS id,
	f.filename AS name,
	f.mimetype,
	f.size::numeric
FROM borehole_file AS bf
JOIN file AS f
	ON f.file_id = bf.file_id
WHERE bf.borehole_id = $1
