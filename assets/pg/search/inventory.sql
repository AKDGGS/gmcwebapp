SELECT
	i.inventory_id AS id,
	co.name AS collection,
	i.sample_number,
	i.slide_number,
	i.box_number,
	i.set_number,
	i.core_number,
	cd.core_diameter,
	i.interval_top,
	i.interval_bottom,
	i.keywords,
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
LEFT OUTER JOIN core_diameter as cd
	on cd.core_diameter_id = i.core_diameter_id
WHERE i.active
GROUP BY i.inventory_id, co.collection_id, cd.core_diameter_id
