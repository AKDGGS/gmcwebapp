SELECT f.file_id AS id, f.description, f.mimetype as type,
	f.size::numeric AS size,
	f.filename AS name
	FROM inventory_file AS ivf
	JOIN file AS f
		ON f.file_id = ivf.file_id
	WHERE ivf.inventory_id = $1
