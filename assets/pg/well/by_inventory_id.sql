SELECT w.well_id, w.name AS well_name, w.alt_names,
	w.well_number, w.api_number,
	w.is_onshore, w.is_federal,
	w.spud_date, w.completion_date,
	w.measured_depth, w.vertical_depth,
	w.elevation, w.elevation_kb,
	w.permit_status, w.completion_status, w.permit_number,
	w.unit::text,
	wo.is_current, o.name AS operator_name, o.remark, ot.name AS operator_type
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
