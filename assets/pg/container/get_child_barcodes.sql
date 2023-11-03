SELECT ARRAY_AGG(barcode::text  ORDER BY barcode::text ) as barcodes,
COUNT(barcode) AS barcode_total
FROM (
WITH RECURSIVE t AS (
	SELECT 0 AS depth, container_id,
		COALESCE(barcode, alt_barcode) AS barcode
	FROM container
	WHERE active
		AND COALESCE(barcode, alt_barcode) = $1

	UNION ALL

	SELECT t.depth + 1 AS depth, c.container_id,
		COALESCE(c.barcode, c.alt_barcode) AS barcode
	FROM container AS c
	JOIN t ON c.parent_container_id = t.container_id
	WHERE active
)
SELECT DISTINCT barcode
FROM (
	SELECT COALESCE(i.barcode, i.alt_barcode) AS barcode
	FROM inventory AS i
	WHERE active
		AND COALESCE(i.barcode, i.alt_barcode) IS NOT NULL
		AND i.container_id IN (SELECT container_id FROM t)

	UNION ALL

	SELECT barcode
	FROM t
) AS q
WHERE barcode != $1
ORDER BY barcode
LIMIT 100
) AS q
