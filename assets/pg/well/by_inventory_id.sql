SELECT w.well_id AS "ID",
	w.name AS "Name",
	w.alt_names "AltNames",
	w.well_number AS "Number",
	w.api_number AS "APINumber",
	w.is_onshore AS "Onshore",
	w.is_federal AS "Federal",
	w.spud_date AS "SpudDate",
	w.completion_date AS "CompletionDate",
	w.measured_depth AS "MeasuredDepth",
	w.vertical_depth AS "VerticalDepth",
	w.elevation,
	w.elevation_kb AS "ElevationKB",
	w.permit_status AS "PermitStatus",
	w.completion_status AS "CompletionStatus",
	w.permit_number AS "PermitNumber",
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
