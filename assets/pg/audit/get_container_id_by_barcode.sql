SELECT container_id
FROM container
WHERE active
	AND barcode = $1 OR alt_barcode = $1
