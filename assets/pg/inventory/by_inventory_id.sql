SELECT
	iv.inventory_id AS id,
	jsonb_build_object(
		'id', cl.collection_id,
		'name', cl.name,
		'description', cl.description
	) AS collection,
	jsonb_build_object(
		'id', co.container_id,
		'name', co.name,
		'path_cache', co.path_cache
	) AS container,
	jsonb_build_object(
		'id', cd.core_diameter_id,
		'name', cd.name,
		'core_diameter', cd.core_diameter,
		'unit', COALESCE(cd.unit::text, 'ft')
	) AS core_diameter,
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
	iv.active
FROM inventory AS iv
LEFT OUTER JOIN collection AS cl
	ON cl.collection_id = iv.collection_id
LEFT OUTER JOIN container AS co
	ON co.container_id = iv.container_id
LEFT OUTER JOIN core_diameter AS cd
	ON cd.core_diameter_id = iv.core_diameter_id
WHERE iv.inventory_id = $1
