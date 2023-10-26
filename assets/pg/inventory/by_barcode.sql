SELECT iv.inventory_id AS id,
	cl.collection_id AS "collection.id",
	cl.name AS "collection.name",
	cl.description AS "collection.description",
	cl.organization_id AS "collection.organization.id",

	co.container_id AS "container.id",
	co.name AS "container.name",
	co.remark AS "container.remark",
	co.barcode AS "container.barcode",
	co.alt_barcode AS "container.alt_barcode",
	co.path_cache AS "container.path_cache",

	iv.dggs_sample_id AS sample_id,
	iv.sample_number,
	iv.sample_number_prefix,
	iv.alt_sample_number,
	iv.published_sample_number,
	iv.published_number_has_suffix,
	iv.barcode,
	iv.alt_barcode,
	iv.state_number,
	iv.box_number,
	iv.set_number,
	iv.split_number,
	iv.slide_number,
	iv.slip_number,
	iv.lab_number,
	iv.lab_report_id,
	iv.map_number,
	iv.description,
	iv.remark,
	iv.tray,
	iv.interval_top,
	iv.interval_bottom,
	iv.keywords::text[],
	COALESCE(iv.interval_unit::text, 'ft') AS interval_unit,
	iv.core_number,

	cd.core_diameter_id AS "core_diameter.id",
	cd.name AS "core_diameter.name",
	cd.core_diameter AS "core_diameter.core_diameter",
	COALESCE(cd.unit::text, 'ft') AS "core_diameter.unit",

	iv.weight,
	iv.weight_unit::text,
	iv.sample_frequency,
	iv.recovery,
	iv.can_publish,
	iv.radiation_msvh,
	iv.received_date,
	iv.entered_date,
	iv.modified_date,
	iv.modified_user,
	iv.active,
	ARRAY_AGG(DISTINCT ib.borehole_id) FILTER (WHERE ib.borehole_id IS NOT NULL) AS borehole_ids,
	ARRAY_AGG(DISTINCT io.outcrop_id) FILTER (WHERE io.outcrop_id IS NOT NULL) AS outcrop_ids,
	ARRAY_AGG(DISTINCT isp.shotpoint_id) FILTER (WHERE isp.shotpoint_id IS NOT NULL) AS shotpoint_ids,
	ARRAY_AGG(DISTINCT iw.well_id) FILTER (WHERE iw.well_id IS NOT NULL) AS well_ids,
	ARRAY_AGG(DISTINCT iq.inventory_quality_id) FILTER (WHERE iq.inventory_quality_id IS NOT NULL) AS quality_ids
FROM inventory AS iv
LEFT OUTER JOIN inventory_borehole AS ib
	ON ib.inventory_id = iv.inventory_id
LEFT OUTER JOIN inventory_outcrop AS io
	ON io.inventory_id = iv.inventory_id
LEFT OUTER JOIN inventory_shotpoint AS isp
	ON isp.inventory_id = iv.inventory_id
LEFT OUTER JOIN inventory_well AS iw
	ON iw.inventory_id = iv.inventory_id
LEFT OUTER JOIN inventory_quality AS iq
	ON iq.inventory_id = iv.inventory_id
LEFT OUTER JOIN collection AS cl
	ON cl.collection_id = iv.collection_id
LEFT OUTER JOIN container AS co
	ON co.container_id = iv.container_id
LEFT OUTER JOIN core_diameter AS cd
	ON cd.core_diameter_id = iv.core_diameter_id
WHERE iv.active AND (iv.barcode = $1
	OR iv.barcode = ('GMC-' || $1)
	OR iv.alt_barcode = $1
	OR iv.container_id IN (
		WITH RECURSIVE r AS (
			SELECT container_id
			FROM container WHERE barcode = $1

			UNION ALL

			SELECT co.container_id
			FROM r
			JOIN container AS co
				ON r.container_id = co.parent_container_id
		) SELECT container_id FROM r
	))
	GROUP BY iv.inventory_id, cl.collection_id, co.container_id, cd.core_diameter_id
LIMIT 100
