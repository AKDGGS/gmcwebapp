SELECT iv.inventory_id AS "Inventory ID",
COALESCE(iv.barcode, iv.alt_barcode, '') AS "Barcode"
FROM inventory AS iv
LEFT OUTER JOIN container AS co
	ON co.container_id = iv.container_id
WHERE iv.active AND iv.container_id IS NULL
