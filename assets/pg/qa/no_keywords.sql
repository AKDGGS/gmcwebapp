SELECT inventory_id,
	COALESCE(co.name || ' ', '')  || COALESCE(c.path_cache, '') AS description
FROM inventory AS i
LEFT OUTER JOIN container AS c
	ON c.container_id = i.container_id AND c.active
LEFT OUTER JOIN collection AS co
	ON co.collection_id = i.collection_id
WHERE i.active AND i.keywords IS NULL
