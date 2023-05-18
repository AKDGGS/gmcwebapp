SELECT well_id AS id,
	name,
	alt_names AS altName,
	well_number AS number,
	api_number AS APINumber,
	is_onshore AS onshore,
	is_federal AS federal,
	permit_status AS permitStatus,
	completion_status AS completionStatus,
	spud_date AS spudDate,
	completion_date AS completionDate,
	measured_depth AS measuredDepth,
	vertical_depth AS verticalDepth,
	elevation,
	elevation_kb AS elevationKb,
	COALESCE(unit::text, 'ft') AS unit
FROM well
WHERE well_id = $1
