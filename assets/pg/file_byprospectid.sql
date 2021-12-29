SELECT f.file_id, f.filename AS file_name,
	pg_size_pretty(f.size::numeric) AS file_size
FROM prospect_file AS pf
JOIN file AS f
	ON f.file_id = pf.file_id
WHERE pf.prospect_id = $1
