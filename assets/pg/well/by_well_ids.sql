SELECT w.well_id AS id,
	w.name AS name,
	w.alt_names,
	w.well_number AS number,
	w.api_number,
	w.is_onshore AS onshore,
	w.is_federal AS federal,
	w.spud_date,
	w.completion_date,
	w.measured_depth,
	w.vertical_depth,
	w.elevation,
	w.elevation_kb,
	w.permit_status,
	w.completion_status,
	w.permit_number,
	w.unit::text,
	ARRAY_AGG(DISTINCT wo.organization_id) FILTER (WHERE wo.organization_id IS NOT NULL) AS operator_ids,
	COALESCE(unit::text, 'ft') AS unit
FROM well AS W
LEFT OUTER JOIN inventory_well AS iw
	ON iw.well_id = w.well_id
LEFT OUTER JOIN well_operator AS wo
	ON wo.well_id = w.well_id
WHERE w.well_id = ANY($1)
GROUP BY w.well_id;
