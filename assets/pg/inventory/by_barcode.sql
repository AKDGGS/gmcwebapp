SELECT iv.inventory_id AS id,
	iv.parent_id AS parentID,

	cl.collection_id AS collectionID,
	cl.name AS "collection.name",
	cl.description AS "collection.description",

	co.container_id AS containerID,
	co.name AS "container.name",
	co.remark AS "container.remark",
	co.barcode AS "container.barcode",
	co.alt_barcode AS "container.altBarcode",
	co.path_cache AS "container.pathCache",

	iv.dggs_sample_id AS sampleID,
	iv.sample_number AS sampleNumber,
	iv.sample_number_prefix AS sampleNumberPrefix,
	iv.alt_sample_number AS altSampleNumber,
	iv.published_sample_number AS publishedSampleNumber,
	iv.published_number_has_suffix AS publishedNumberHasSuffix,
	iv.barcode AS barcode,
	iv.alt_barcode AS altBarcode,
	iv.state_number AS stateNumber,
	iv.box_number AS boxNumber,
	iv.set_number AS setNumber,
	iv.split_number AS splitNumber,
	iv.slide_number AS slideNumber,
	iv.slip_number AS slipNumber,
	iv.lab_number AS labNumber,
	iv.lab_report_id AS labReportID,
	iv.map_number AS mapNumber,
	iv.description AS description,
	iv.remark AS remark,
	iv.tray AS tray,
	iv.interval_top AS intervalTop,
	iv.interval_bottom AS intervalBottom,
	iv.keywords::text[] AS keywords,
	COALESCE(iv.interval_unit::text, 'ft') AS intervalUnit,
	iv.core_number AS coreNumber,

	cd.core_diameter_id AS coreDiameterID,
	cd.name AS "coreDiameter.name",
	cd.core_diameter AS "coreDiameter.coreDiameter",
	COALESCE(cd.unit::text, 'ft') AS "coreDiameter.Unit",

	iv.weight AS weight,
	iv.weight_unit::text AS weightUnit,
	iv.sample_frequency AS sampleFrequency,
	iv.recovery AS recovery,
	iv.can_publish AS canPublish,
	iv.radiation_msvh AS radiationMSVH,
	iv.received_date AS receivedDate,
	iv.entered_date AS enteredDate,
	iv.modified_date AS modifiedDate,
	iv.modified_user AS modifiedUser,
	iv.active AS active
FROM inventory AS iv
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
LIMIT 100
