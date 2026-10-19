SELECT
	f.file_id AS id,
	f.filename AS name,
	f.mimetype,
	f.size::numeric,
	f.file_type,
	array_position(enum_range(NULL::file_types), file_type) AS file_type_int
FROM borehole_file AS bf
JOIN file AS f
	ON f.file_id = bf.file_id
WHERE bf.borehole_id = $1
ORDER BY f.file_type NULLS LAST;
