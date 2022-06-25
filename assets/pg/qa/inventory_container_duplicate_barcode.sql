SELECT i.inventory_id AS "Inventory ID",
	c.container_id AS "Container ID",
	i.barcode AS "Inv Barcode",
	i.alt_barcode AS "Inv Alt Barcode",
	c.barcode AS "Con Barcode",
	c.alt_barcode AS "Con Alt Barcode"
FROM inventory AS i
JOIN container AS c ON
	(
		i.barcode IS NOT NULL AND (
			i.barcode = c.barcode OR i.barcode = c.alt_barcode
		)
	) OR (
		i.alt_barcode IS NOT NULL AND (
			i.alt_barcode = c.barcode OR i.alt_barcode = c.alt_barcode
		)
	)
WHERE i.active AND c.active
