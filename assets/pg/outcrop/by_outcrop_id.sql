SELECT
	oc.outcrop_id AS id,
	oc.name,
	oc.outcrop_number AS "number",
	oc.is_onshore AS onshore,
	oc.year,
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
			JOIN outcrop_organization AS oo
				ON o.organization_id = oo.organization_id
			WHERE oo.outcrop_id = oc.outcrop_id
			ORDER BY o.name
		) AS s
	) AS organizations
FROM outcrop AS oc
WHERE oc.outcrop_id = $1
