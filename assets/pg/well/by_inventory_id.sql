SELECT w.well_id, w.name AS "WellName", w.alt_names "AltNames",
	w.well_number AS "WellNumber",
	w.api_number AS "APINumber",
	w.is_onshore AS "IsOnshore", 
	w.is_federal AS "IsFederal",
	w.spud_date AS "SpudDate", 
	w.completion_date AS "CompletionDate",
	w.measured_depth AS "MeasuredDepth",
	w.vertical_depth AS "VerticalDepth",
	w.elevation,
	w.elevation_kb AS "ElevationKB",
	w.permit_status AS "PermitStatus",
	w.completion_status AS "ComletionStatus",
	w.permit_number AS "PermitNumber",
	w.unit::text,
	wo.is_current AS "IsCurrent",
	o.name AS "OperatorName",
	o.remark,
	ot.name AS "OperatorType"
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
