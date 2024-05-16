SELECT
	f.file_id AS id,
	f.filename AS name,
	f.mimetype,
	f.size::numeric
FROM well_file AS wf
JOIN file AS f
	ON f.file_id = wf.file_id
WHERE wf.well_id = $1
