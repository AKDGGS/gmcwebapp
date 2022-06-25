SELECT i.inventory_id AS "Inventory ID",
	c.path_cache AS "Container",
	i.remark AS "Remark"
FROM inventory AS i
LEFT OUTER JOIN container AS c
	ON c.container_id = i.container_id
WHERE i.active AND 
	COALESCE(i.barcode, i.alt_barcode, c.barcode, c.alt_barcode) IS NULL
