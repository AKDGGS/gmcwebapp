SELECT f.file_id AS "ID", f.filename AS "Name", f.mimetype as "Type",
	f.size::numeric AS "Size"
FROM outcrop_file AS of
JOIN file AS f
	ON f.file_id = of.file_id
WHERE of.outcrop_id = $1
