SELECT
	f.file_id AS id,
	f.filename AS name,
	f.mimetype,
	f.size::numeric AS "size"
FROM outcrop_file AS ocf
JOIN file AS f
	ON f.file_id = ocf.file_id
WHERE ocf.outcrop_id = $1
