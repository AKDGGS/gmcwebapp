SELECT
	i.inventory_id AS id,
	co.name AS collection,
	i.barcode,
	i.remark,
	i.can_publish,
	jsonb_agg(ST_ASGeoJSON(ig.geog)) FILTER (
		WHERE ig.geog IS NOT NULL AND ST_NPoints(ig.geog::geometry) < 100
	) AS geometries
FROM inventory AS i
LEFT OUTER JOIN inventory_geog AS ig
	ON ig.inventory_id = i.inventory_id
LEFT OUTER JOIN collection AS co
	ON co.collection_id = i.collection_id
WHERE i.active
GROUP BY i.inventory_id, co.collection_id
