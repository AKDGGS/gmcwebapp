SELECT io.inventory_id AS "Inventory ID",
	COALESCE(iv.barcode, iv.alt_barcode) AS "Barcode",
	ct.path_cache AS "Container"
FROM inventory_outcrop AS io
JOIN inventory AS iv ON iv.inventory_id = io.inventory_id
LEFT OUTER JOIN container AS ct
	ON ct.container_id = iv.container_id
WHERE iv.active AND (
	iv.sample_number IS NULL OR LENGTH(TRIM(BOTH FROM iv.sample_number)) = 0
)
ORDER BY io.inventory_id ASC
