SELECT SUM(count) AS count
FROM (
	SELECT COUNT(container_id) AS count
	FROM container
	WHERE COALESCE(barcode, alt_barcode) = $1

	UNION ALL

	SELECT COUNT(inventory_id) AS count
	FROM inventory
	WHERE COALESCE(barcode, alt_barcode) = $1
) AS q
