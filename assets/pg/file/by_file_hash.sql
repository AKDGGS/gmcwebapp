SELECT
	f.file_id AS id,
	f.filename AS name,
	f.file_hash,
	f.description,
	f.mimetype,
	f.size::numeric,
	COALESCE(f.file_type::text, '') AS file_type,
	array_position(enum_range(NULL::file_types), file_type) AS file_type_int
FROM file AS f
WHERE f.file_hash = $1
