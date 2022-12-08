SELECT f.file_id AS ID, f.filename AS Name,
	pg_size_pretty(f.size::numeric) AS "Size"
FROM outcrop_file AS of
JOIN file AS f
	ON f.file_id = of.file_id
WHERE of.outcrop_id = $1
