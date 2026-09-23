SELECT
	f.file_id AS id,
	f.filename AS name,
	f.mimetype,
	f.size::numeric,
	COALESCE(f.file_type::text, '') AS file_type,
	array_position(enum_range(NULL::file_types), file_type) AS file_type_int
FROM outcrop_file AS ocf
JOIN file AS f
	ON f.file_id = ocf.file_id
WHERE ocf.outcrop_id = $1
ORDER BY f.file_type NULLS LAST;