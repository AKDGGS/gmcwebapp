SELECT
	f.file_id AS id,
	f.filename AS name,
	f.mimetype,
	f.size::numeric
FROM prospect_file AS pf
JOIN file AS f
	ON f.file_id = pf.file_id
WHERE pf.prospect_id = $1
