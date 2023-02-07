SELECT p.prospect_id AS "Prospect.ID",
	p.name AS "Prospect.Name",
	p.alt_names AS "Prospect.AltProspectNames", 
	p.ardf_number AS "Prospect.ARDFNumber",
	b.borehole_id AS "ID", 
	b.name AS "Name", 
	b.alt_names AS "AltNames",
	b.is_onshore AS "IsOnshore", 
	b.completion_date AS "CompletionDate",
	b.measured_depth AS "MeasuredDepth",
	COALESCE(measured_depth_unit::text, 'ft') AS "MeasuredDepthUnit",
	b.elevation AS "Elevation",
	COALESCE(elevation_unit::text, 'ft') AS "ElevationUnit"
FROM borehole AS b
LEFT OUTER JOIN prospect AS p
	ON p.prospect_id = b.prospect_id
WHERE b.borehole_id = $1
