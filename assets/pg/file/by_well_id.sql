SELECT
	f.file_id AS id,
	f.filename AS name,
	f.mimetype,
	f.size::numeric,
	f.file_type,
	array_position(enum_range(NULL::file_types), file_type) AS file_type_int
FROM well_file AS wf
JOIN file AS f
	ON f.file_id = wf.file_id
WHERE wf.well_id = $1
ORDER BY f.file_type NULLS LAST;
