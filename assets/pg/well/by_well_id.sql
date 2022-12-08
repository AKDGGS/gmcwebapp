SELECT well_id AS ID, 
	name, 
	alt_names AS AltName, 
	well_number AS WellNumber,
	api_number AS APINumber,
	is_onshore AS IsOnshore,
	is_federal AS IsFederal, 
	permit_status AS PermitStatus, 
	completion_status AS CompletionStatus, 
	spud_date AS SpudDate,
	completion_date AS CompletionDate,
	measured_depth AS MeasuredDepth,
	vertical_depth AS VerticalDepth,
	elevation,
	elevation_kb AS ElevationKb,
	COALESCE(unit::text, 'ft') AS Unit
FROM well
WHERE well_id = $1
