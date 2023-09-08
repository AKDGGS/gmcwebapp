SELECT well_id AS id,
	name,
	alt_names,
	well_number AS number,
	api_number,
	is_onshore AS onshore,
	is_federal AS federal,
	permit_status,
	permit_number,
	completion_status,
	spud_date,
	completion_date,
	measured_depth,
	vertical_depth,
	elevation,
	elevation_kb,
	COALESCE(unit::text, 'ft') AS unit
FROM well
WHERE well_id = $1
