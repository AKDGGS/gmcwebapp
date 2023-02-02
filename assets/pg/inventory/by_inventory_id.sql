SELECT iv.inventory_id AS "ID",
	iv.parent_id AS "ParentID",

	cl.collection_id AS "CollectionID",
	cl.name AS "CollectionName",
	cl.description AS "CollectionDescription",

	co.container_id AS "ContainerID",
	co.name AS "ContainerName",
	co.remark AS "ContainerRemark",
	co.barcode AS "ContainerBarcode",
	co.alt_barcode AS "ContainerAltBarcode",
	co.path_cache AS "ContainerPathCache",

	iv.dggs_sample_id AS "SampleID",
	iv.sample_number AS "SampleNumber",
	iv.sample_number_prefix AS "SampleNumberPrefix",
	iv.alt_sample_number AS "AltSampleNumber",
	iv.published_sample_number AS "PublishedSampleNumber",
	iv.published_number_has_suffix AS "PublishedNumberHasSuffix",
	iv.barcode AS "Barcode",
	iv.alt_barcode AS "AltBarcode",
	iv.state_number AS "StateNumber",
	iv.box_number AS "BoxNumber",
	iv.set_number AS "SetNumber",
	iv.split_number AS "SplitNumber",
	iv.slide_number AS "SlideNumber",
	iv.slip_number AS "SlipNumber",
	iv.lab_number AS "LabNumber",
	iv.lab_report_id AS "LabReportID",
	iv.map_number AS "MapNumber",
	iv.description AS "Description",
	iv.remark AS "Remark",
	iv.tray AS "Tray",
	iv.interval_top AS "IntervalTop",
	iv.interval_bottom AS "IntervalBottom",
	iv.keywords::text[] AS "Keywords",
	COALESCE(iv.interval_unit::text, 'ft') AS IntervalUnit,
	iv.core_number AS "CoreNumber",
	
	cd.core_diameter_id AS "CoreDiameterID",
	cd.name AS "CoreDiameterName",
	cd.core_diameter AS "CoreDiameter",
	COALESCE(cd.unit::text, 'ft') AS "CoreDiameterUnit",
	
	iv.weight AS "Weight",
	iv.weight_unit::text AS "WeightUnit",
	iv.sample_frequency AS "SampleFrequency",
	iv.recovery AS "Recovery",
	iv.can_publish AS "CanPublish",
	iv.radiation_msvh AS "RadiationMSVH",
	iv.received_date AS "ReceivedDate",
	iv.entered_date AS "EnteredDate",
	iv.modified_date AS "ModifiedDate",
	iv.modified_user AS "ModifiedUser",
	iv.active AS "Active"
FROM inventory AS iv
LEFT OUTER JOIN collection AS cl
	ON cl.collection_id = iv.collection_id
LEFT OUTER JOIN container AS co
	ON co.container_id = iv.container_id
LEFT OUTER JOIN core_diameter AS cd
	ON cd.core_diameter_id = iv.core_diameter_id
WHERE iv.inventory_id = $1
