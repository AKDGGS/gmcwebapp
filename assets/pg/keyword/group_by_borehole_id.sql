SELECT keywords AS keywords,
	COUNT(q2.inventory_id) AS count
FROM (
	SELECT q1.inventory_id,
		ARRAY_AGG(q1.keyword ORDER BY q1.keyword)::text[] AS keywords
	FROM (
		SELECT i.inventory_id, UNNEST(i.keywords) AS keyword
		FROM inventory_borehole AS ib
		JOIN borehole AS b
			ON b.borehole_id = ib.borehole_id
		JOIN inventory AS i
			ON i.inventory_id = ib.inventory_id
		WHERE b.borehole_id = $1 AND i.active
			AND (i.can_publish = true OR i.can_publish = $2)
		GROUP BY i.inventory_id
	) AS q1
	GROUP BY q1.inventory_id
) AS q2
GROUP BY q2.keywords
ORDER BY q2.keywords
