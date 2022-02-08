SELECT ARRAY_TO_JSON(keywords) AS keywords,
	COUNT(q2.inventory_id) AS count
FROM (
	SELECT q1.inventory_id,
		ARRAY_AGG(q1.keyword ORDER BY q1.keyword) AS keywords
	FROM (
		SELECT i.inventory_id, UNNEST(i.keywords) AS keyword
		FROM inventory_outcrop AS io
		JOIN outcrop AS o
			ON o.outcrop_id = io.outcrop_id
		JOIN inventory AS i
			ON i.inventory_id = io.inventory_id
		WHERE o.outcrop_id = $1 AND i.active
			AND (i.can_publish = true OR i.can_publish = $2)
		GROUP BY i.inventory_id
	) AS q1
	GROUP BY q1.inventory_id
) AS q2
GROUP BY q2.keywords
ORDER BY q2.keywords
