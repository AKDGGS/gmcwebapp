SELECT f.file_id, f.filename AS file_name,
	pg_size_pretty(f.size::numeric) AS file_size
FROM borehole_file AS bf
JOIN file AS f
	ON f.file_id = bf.file_id
WHERE bf.borehole_id = $1
