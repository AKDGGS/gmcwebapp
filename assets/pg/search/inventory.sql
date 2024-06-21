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
	w.wells,
	g.geometries
FROM inventory AS i
LEFT OUTER JOIN (
	SELECT
		iw.inventory_id,
		jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
			'id', w.well_id,
			'name', w.name,
			'altnames', w.alt_names,
			'number', w.well_number,
			'api', w.api_number
		))) AS wells
	FROM inventory_well AS iw
	JOIN well AS w ON w.well_id = iw.well_id
	GROUP BY iw.inventory_id
) AS w ON w.inventory_id = i.inventory_id
LEFT OUTER JOIN (
	SELECT
		inventory_id,
		jsonb_agg(ST_ASGeoJSON(geog)::jsonb) AS geometries
	FROM inventory_geog
	WHERE ST_NPoints(geog::geometry) < 100
	GROUP BY inventory_id
) AS g ON g.inventory_id = i.inventory_id
LEFT OUTER JOIN collection AS co
	ON co.collection_id = i.collection_id
LEFT OUTER JOIN core_diameter as cd
	on cd.core_diameter_id = i.core_diameter_id
WHERE i.active
