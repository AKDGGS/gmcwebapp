UPDATE container AS c1
	SET parent_container_id = $2
FROM container AS c2
WHERE c1.parent_container_id = c2.container_id
	AND (
		c2.barcode = $1
		OR c2.barcode = ('GMC-' || $1)
		OR c2.alt_barcode = $1
	)
