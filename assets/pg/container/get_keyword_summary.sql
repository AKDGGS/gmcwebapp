SELECT ARRAY_AGG(keyword::text ORDER BY keyword::text) as keywords
FROM (
	SELECT DISTINCT UNNEST(keywords) AS keyword
	FROM inventory
	WHERE active AND container_id IN (
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
	LIMIT 100
) AS q
