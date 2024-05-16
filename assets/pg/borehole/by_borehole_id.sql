SELECT
	b.borehole_id AS id,
	b.name,
	b.alt_names,
	b.is_onshore AS onshore,
	b.completion_date,
	b.measured_depth,
	COALESCE(measured_depth_unit::text, 'ft') AS measured_depth_unit,
	b.elevation,
	COALESCE(elevation_unit::text, 'ft') AS elevation_unit,
	jsonb_build_object(
		'id', p.prospect_id,
		'name', p.name,
		'alt_names', p.alt_names,
		'ardf', p.ardf_number
	) AS prospect,
	(
		SELECT jsonb_agg(orgs)
		FROM (
			SELECT jsonb_build_object(
				'id', o.organization_id,
				'name', o.name,
				'type', jsonb_build_object(
					'name', ot.name
				),
				'remark', o.remark
			) AS orgs
			FROM organization AS o
			JOIN organization_type AS ot
				ON o.organization_type_id = ot.organization_type_id
			JOIN borehole_organization AS bo
				ON o.organization_id = bo.organization_id
			WHERE bo.borehole_id = b.borehole_id
			ORDER BY o.name
		) AS s
	) AS organizations
FROM borehole AS b
LEFT OUTER JOIN prospect AS p
	ON p.prospect_id = b.prospect_id
WHERE b.borehole_id = $1
