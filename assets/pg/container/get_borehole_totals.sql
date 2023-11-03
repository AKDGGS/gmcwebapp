SELECT prospect, borehole, borehole_total
FROM (
	SELECT ps.prospect_id, bh.borehole_id,
		ps.name AS prospect, bh.name AS borehole,
		COUNT(DISTINCT COALESCE(iv.barcode, iv.alt_barcode)) AS borehole_total
	FROM container AS co
	JOIN inventory AS iv
		ON iv.container_id = co.container_id
	JOIN inventory_borehole AS ivb
		ON ivb.inventory_id = iv.inventory_id
	JOIN borehole AS bh
		ON bh.borehole_id = ivb.borehole_id
	LEFT OUTER JOIN prospect AS ps
		ON ps.prospect_id = bh.prospect_id
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
	GROUP BY ps.prospect_id, bh.borehole_id, ps.name, bh.name
	LIMIT 100
) AS q
ORDER BY prospect, borehole
