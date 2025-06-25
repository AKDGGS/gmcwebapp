SELECT barcode,
	alt_barcode
FROM container
WHERE active
	AND parent_container_id = $1
	AND (barcode IS NOT NULL OR alt_barcode IS NOT NULL)

UNION

SELECT barcode,
	alt_barcode
FROM inventory
WHERE active
	AND container_id = $1
	AND (barcode IS NOT NULL OR alt_barcode IS NOT NULL)
