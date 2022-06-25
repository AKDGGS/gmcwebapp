SELECT container_id AS "Container ID",
	path_cache AS "Container"
FROM container
WHERE active AND (
	(barcode IS NOT NULL AND LENGTH(TRIM(BOTH FROM barcode)) = 0)
	OR
	(alt_barcode IS NOT NULL AND LENGTH(TRIM(BOTH FROM alt_barcode)) = 0)
)
