SELECT
	w.well_id AS id,
	w.name,
	w.alt_names,
	w.well_number AS "number",
	w.api_number,
	w.is_onshore AS onshore,
	w.is_federal AS federal,
	w.permit_status,
	w.permit_number,
	w.completion_status,
	w.spud_date,
	w.completion_date,
	w.measured_depth,
	w.vertical_depth,
	w.elevation,
	w.elevation_kb,
	COALESCE(unit::text, 'ft') AS unit,
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
WHERE w.well_id = $1
