SELECT barcode AS "Barcode",
	STRING_AGG(path_cache, ', ') AS "Containers"
FROM (
	SELECT DISTINCT COALESCE(iv.barcode, iv.alt_barcode) AS barcode,
		iv.container_id, co.path_cache
	FROM inventory AS iv
	JOIN container AS co
		ON co.container_id = iv.container_id
	WHERE iv.active AND (iv.barcode IS NOT NULL
		OR iv.alt_barcode IS NOT NULL)
) AS q1
GROUP BY barcode
HAVING COUNT(DISTINCT container_id) > 1

