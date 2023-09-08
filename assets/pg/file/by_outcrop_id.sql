SELECT f.file_id AS id,
	f.filename AS name,
	f.mimetype AS type,
	f.size::numeric AS size
FROM outcrop_file AS of
JOIN file AS f
	ON f.file_id = of.file_id
WHERE of.outcrop_id = $1
