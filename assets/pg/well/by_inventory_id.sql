SELECT w.well_id AS id,
	w.name AS name,
	w.alt_names altNames,
	w.well_number AS number,
	w.api_number AS APINumber,
	w.is_onshore AS onshore,
	w.is_federal AS federal,
	w.spud_date AS spudDate,
	w.completion_date AS completionDate,
	w.measured_depth AS measuredDepth,
	w.vertical_depth AS verticalDepth,
	w.elevation,
	w.elevation_kb AS elevationKB,
	w.permit_status AS permitStatus,
	w.completion_status AS completionStatus,
	w.permit_number AS permitNumber,
	w.unit::text
	FROM well AS w
	LEFT OUTER JOIN inventory_well AS iw
		ON iw.well_id = w.well_id
	WHERE iw.inventory_id = $1
	GROUP BY w.well_id;
