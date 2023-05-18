SELECT bh.borehole_id AS id,
	bh.name,
	bh.alt_names AS altNames,
	bh.is_onshore AS onshore,
	bh.completion_date AS completionDate,
	bh.measured_depth AS measuredDepth,
	COALESCE(bh.measured_depth_unit::text, 'ft') AS measuredDepthUnit,
	bh.elevation AS elevation,
	COALESCE(bh.elevation_unit::text, 'ft') AS elevationUnit,
	ph.prospect_id AS "prospect.id",
	ph.name AS "prospect.name",
	ph.alt_names AS "prospect.altNames",
	ph.ardf_number AS "prospect.ARDFNumber"
FROM borehole AS bh
LEFT OUTER JOIN prospect AS ph
	ON ph.prospect_id = bh.prospect_id
JOIN inventory_borehole AS ib
	ON ib.borehole_id = bh.borehole_id
WHERE ib.inventory_id = $1
