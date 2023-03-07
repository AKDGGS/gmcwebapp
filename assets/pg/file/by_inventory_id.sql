SELECT f.file_id AS "ID", f.description, f.mimetype as "Type",
	pg_size_pretty(f.size::numeric) AS "Size",
	f.filename AS "Name"
	FROM inventory_file AS ivf
	JOIN file AS f
		ON f.file_id = ivf.file_id
	WHERE ivf.inventory_id = $1
