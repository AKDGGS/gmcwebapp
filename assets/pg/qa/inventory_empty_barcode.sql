SELECT inventory_id AS "Inventory ID",
	remark AS "Remark"
FROM inventory
WHERE active AND (
	(barcode IS NOT NULL AND LENGTH(TRIM(BOTH FROM barcode)) = 0)
	OR
	(alt_barcode IS NOT NULL AND LENGTH(TRIM(BOTH FROM alt_barcode)) = 0)
)
