SELECT i.inventory_id, co.name AS collection,
	i.barcode, i.remark,
	jsonb_agg(ST_ASGeoJSON(ig.geog)) FILTER (
		WHERE ig.geog IS NOT NULL
	) AS geometries
FROM inventory AS i
LEFT OUTER JOIN inventory_geog AS ig
	ON ig.inventory_id = i.inventory_id
LEFT OUTER JOIN collection AS co
	ON co.collection_id = i.collection_id
WHERE ST_NPoints(ig.geog::geometry) < 100
GROUP BY i.inventory_id, co.collection_id
