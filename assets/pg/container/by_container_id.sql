SELECT container_id AS id,
	name,
	path_cache,
	remark,
	barcode,
	alt_barcode
FROM container
WHERE container_id = $1
