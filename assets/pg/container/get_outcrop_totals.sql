SELECT outcrop, outcrop_total
		FROM (
			SELECT oc.outcrop_id,
				(oc.name || COALESCE(' - ' || oc.outcrop_number, '')) AS outcrop,
				COUNT(DISTINCT COALESCE(iv.barcode, iv.alt_barcode)) AS outcrop_total
			FROM container AS co
			JOIN inventory AS iv ON iv.container_id = co.container_id
			JOIN inventory_outcrop AS ivo ON ivo.inventory_id = iv.inventory_id
			JOIN outcrop AS oc ON oc.outcrop_id = ivo.outcrop_id
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
			GROUP BY oc.outcrop_id, oc.name
			LIMIT 100
		) AS q
		ORDER BY outcrop
