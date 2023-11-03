SELECT well, well_total
FROM (
	SELECT we.well_id,
		(we.name || COALESCE(' - ' || we.well_number, '')) AS well,
		COUNT(DISTINCT COALESCE(iv.barcode, iv.alt_barcode)) AS well_total
	FROM container AS co
	JOIN inventory AS iv ON iv.container_id = co.container_id
	JOIN inventory_well AS ivw ON ivw.inventory_id = iv.inventory_id
	JOIN well AS we ON we.well_id = ivw.well_id
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
	GROUP BY we.well_id, we.name
	LIMIT 100
) AS q
ORDER BY well
