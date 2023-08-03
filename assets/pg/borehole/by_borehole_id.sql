SELECT p.prospect_id AS "prospect.id",
	p.name AS "prospect.name",
	p.alt_names AS "prospect.altNames",
	p.ardf_number AS "prospect.ARDFNumber",
	b.borehole_id AS id,
	b.name,
	b.alt_names AS altNames,
	b.is_onshore AS onshore,
	b.completion_date AS completionDate,
	b.measured_depth AS measuredDepth,
	COALESCE(measured_depth_unit::text, 'ft') AS measuredDepthUnit,
	b.elevation AS elevation,
	COALESCE(elevation_unit::text, 'ft') AS elevationUnit
FROM borehole AS b
LEFT OUTER JOIN prospect AS p
	ON p.prospect_id = b.prospect_id
WHERE b.borehole_id = $1
