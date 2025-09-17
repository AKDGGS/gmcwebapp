SELECT
	(
		SELECT COUNT(*)
		FROM inventory
		WHERE active
			AND (barcode = $1 OR alt_barcode = $1
	) AS inventory_count,
	(
		SELECT COUNT(*)
		FROM container
		WHERE active
			AND (barcode = $1 OR alt_barcode = $1)
	) AS container_count
