SELECT
	f.file_id AS id,
	f.filename AS name,
	f.description,
	f.mimetype,
	f.size::numeric,
	f.file_type,
	array_position(enum_range(NULL::file_types), file_type) AS file_type_int
FROM file AS f
WHERE f.file_id = $1
