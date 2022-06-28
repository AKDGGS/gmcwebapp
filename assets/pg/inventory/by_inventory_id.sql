SELECT iv.inventory_id,
	iv.parent_id,

	cl.collection_id,
	cl.name AS collection_name,
	cl.description AS collection_description,

	co.container_id,
	co.name AS container_name,
	co.remark AS container_remark,
	co.barcode AS container_barcode,
	co.alt_barcode AS container_alt_barcode,
	co.path_cache AS container_path_cache,

	iv.dggs_sample_id,
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
	ARRAY_TO_JSON(iv.keywords) AS keywords,
	COALESCE(iv.interval_unit::text, 'ft') AS interval_unit,
	iv.core_number,

	cd.core_diameter_id,
	cd.name AS core_diameter_name,
	cd.core_diameter,
	COALESCE(cd.unit::text, 'ft') AS unit,
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
