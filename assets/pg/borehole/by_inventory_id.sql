SELECT
	bh.borehole_id AS id,
	bh.name,
	bh.alt_names,
	bh.is_onshore,
	bh.completion_date,
	bh.measured_depth,
	COALESCE(bh.measured_depth_unit::text, 'ft') AS measured_depth_unit,
	bh.elevation,
	COALESCE(bh.elevation_unit::text, 'ft') AS elevation_unit,
	bh.stash,
	bh.entered_date,
	bh.modified_date,
	bh.modified_user,
	jsonb_build_object(
		'id', ph.prospect_id,
		'name', ph.name,
		'ardf_number', ph.ardf_number
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
			WHERE bo.borehole_id = bh.borehole_id
			ORDER BY o.name
		) AS s
	) AS organizations
FROM borehole AS bh
LEFT OUTER JOIN prospect AS ph
	ON ph.prospect_id = bh.prospect_id
JOIN inventory_borehole AS ib
	ON ib.borehole_id = bh.borehole_id
LEFT JOIN borehole_point AS bp
	ON bp.borehole_id = bh.borehole_id
LEFT JOIN point AS p
	ON p.point_id = bp.point_id
LEFT JOIN mining_district AS md
	ON ST_Intersects(md.geog, p.geog)
WHERE ib.inventory_id = $1
