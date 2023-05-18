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
	w.unit::text,
	json_agg(json_build_object('Name', o.name, 'Remark', o.remark, 'Type', ot.name, 'Current', wo.is_current) ORDER BY wo.is_current DESC) AS "Organizations"
	FROM well AS w
	LEFT OUTER JOIN well_operator AS wo
		ON wo.well_id = w.well_id
	LEFT OUTER JOIN organization AS o
		ON wo.organization_id = o.organization_id
	LEFT OUTER JOIN organization_type AS ot
		ON o.organization_type_id= ot.organization_type_id
	LEFT OUTER JOIN inventory_well AS iw
		ON iw.well_id = w.well_id
	WHERE iw.inventory_id = $1
	GROUP BY w.well_id;
