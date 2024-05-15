SELECT
	f.file_id AS id,
	f.filename AS name,
	f.description,
	f.mimetype,
	f.size::numeric
FROM inventory_file AS ivf
JOIN file AS f
	ON f.file_id = ivf.file_id
WHERE ivf.inventory_id = $1
