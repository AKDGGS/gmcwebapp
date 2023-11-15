UPDATE inventory AS i
	SET container_id = $2
FROM container AS c
WHERE c.container_id = i.container_id
	AND (
		c.barcode = $1
		OR c.barcode = ('GMC-' || $1)
		OR c.alt_barcode = $1
	)
