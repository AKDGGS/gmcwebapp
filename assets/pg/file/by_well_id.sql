SELECT f.file_id, f.filename AS file_name,
	pg_size_pretty(f.size::numeric) AS file_size
FROM well_file AS wf
JOIN file AS f
	ON f.file_id = wf.file_id
WHERE wf.well_id = $1
