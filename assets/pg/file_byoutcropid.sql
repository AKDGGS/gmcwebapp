SELECT f.file_id, f.filename AS file_name,
	pg_size_pretty(f.size::numeric) AS file_size
FROM outcrop_file AS of
JOIN file AS f
	ON f.file_id = of.file_id
WHERE of.outcrop_id = $1
