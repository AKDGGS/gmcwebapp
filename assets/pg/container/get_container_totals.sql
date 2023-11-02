SELECT path_cache, total AS container_total
FROM (
	SELECT c.container_id, c.path_cache, COUNT(i.inventory_id) AS total
	FROM inventory AS i
	JOIN container AS c ON c.container_id = i.container_id
	WHERE i.active AND i.container_id IN (
		WITH RECURSIVE t AS (
			SELECT 0 AS depth, container_id
			FROM container
			WHERE active AND COALESCE(barcode, alt_barcode) = $1

			UNION ALL

			SELECT t.depth + 1 AS depth, c.container_id
			FROM container AS c
			JOIN t ON c.parent_container_id = t.container_id
			WHERE active
		) SELECT container_id FROM t
	)
	GROUP BY c.container_id, c.path_cache
	LIMIT 100
) AS q
ORDER BY path_cache
