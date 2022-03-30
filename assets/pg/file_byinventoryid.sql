SELECT f.file_id, f.description, f.mimetype,
	pg_size_pretty(f.size::numeric) AS file_size,
	f.filename AS file_name,
	ft.file_type_id, ft.name
	FROM inventory_file AS ivf
	JOIN file AS f
		ON f.file_id = ivf.file_id
	LEFT OUTER JOIN file_type AS ft
		ON ft.file_type_id = f.file_type_id
	WHERE ivf.inventory_id = $1
