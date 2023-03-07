SELECT bh.borehole_id AS "ID",
	bh.name,
	bh.alt_names AS "AltNames",
	bh.is_onshore AS "Onshore",
	bh.completion_date AS "CompletionDate",
	bh.measured_depth AS "MeasuredDepth",
	COALESCE(bh.measured_depth_unit::text, 'ft') AS "MeasuredDepthUnit",
	bh.elevation AS "Elevation",
	COALESCE(bh.elevation_unit::text, 'ft') AS "ElevationUnit",
	ph.prospect_id AS "Prospect.ID",
	ph.name AS "Prospect.Name",
	ph.alt_names AS "Prospect.AltNames",
	ph.ardf_number AS "Prospect.ARDFNumber"
FROM borehole AS bh
LEFT OUTER JOIN prospect AS ph
	ON ph.prospect_id = bh.prospect_id
JOIN inventory_borehole AS ib
	ON ib.borehole_id = bh.borehole_id
WHERE ib.inventory_id = $1
