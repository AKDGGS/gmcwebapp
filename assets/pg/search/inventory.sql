SELECT i.inventory_id, i.barcode,
	jsonb_agg(ST_ASGeoJSON(ig.geog)) FILTER (WHERE ig.geog IS NOT NULL)
	AS geometries
FROM inventory AS i
LEFT OUTER JOIN inventory_geog AS ig
	ON ig.inventory_id = i.inventory_id
GROUP BY i.inventory_id
