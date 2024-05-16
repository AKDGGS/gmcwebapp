SELECT w.well_id AS id,
	w.name AS name,
	w.alt_names,
	w.well_number AS "number",
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
	(
		SELECT jsonb_agg(operators)
		FROM (
			SELECT jsonb_build_object(
				'organization_id', o.organization_id,
				'name', o.name,
				'type', jsonb_build_object(
					'name', ot.name
				),
				'remark', o.remark,
				'is_current', wo.is_current
			) AS operators
			FROM organization AS o
			JOIN organization_type AS ot
				ON o.organization_type_id = ot.organization_type_id
			JOIN well_operator AS wo
				ON o.organization_id = wo.organization_id
			WHERE wo.well_id = w.well_id
			ORDER BY wo.is_current DESC, o.name
		) AS s
	) AS organizations
FROM well AS w
LEFT OUTER JOIN inventory_well AS iw
	ON iw.well_id = w.well_id
WHERE iw.inventory_id = $1
GROUP BY w.well_id
