SELECT collection, collection_total
FROM (
	SELECT co.collection_id, co.name AS collection,
		COUNT(DISTINCT COALESCE(iv.barcode, iv.alt_barcode)) AS collection_total
	FROM inventory AS iv
	JOIN collection AS co ON co.collection_id = iv.collection_id
	WHERE iv.active AND iv.container_id IN (
		WITH RECURSIVE t AS (
			SELECT 0 AS depth, container_id
			FROM container
			WHERE active
				AND COALESCE(barcode, alt_barcode) = $1

			UNION ALL

			SELECT t.depth + 1 AS depth, c.container_id
			FROM container AS c
			JOIN t ON c.parent_container_id = t.container_id
			WHERE active
		) SELECT container_id FROM t
	)
	GROUP BY co.collection_id, co.name
	LIMIT 100
) AS q
ORDER BY collection
