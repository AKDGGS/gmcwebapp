SELECT shotline, total
FROM (
	SELECT sl.shotline_id, sl.name AS shotline,
		COUNT(DISTINCT COALESCE(iv.barcode, iv.alt_barcode)) AS total
	FROM inventory AS iv
	JOIN inventory_shotpoint AS ivs
		ON ivs.inventory_id = iv.inventory_id
	JOIN shotpoint AS sp
		ON sp.shotpoint_id = ivs.shotpoint_id
	JOIN shotline AS sl
		ON sl.shotline_id = sp.shotline_id
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
	GROUP BY sl.shotline_id, sl.name
	LIMIT 100
) AS q
ORDER BY shotline
