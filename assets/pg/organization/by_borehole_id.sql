SELECT o.organization_id AS id, o.name, ot.name AS type, o.remark
FROM organization AS o
JOIN organization_type AS ot
	ON o.organization_type_id = ot.organization_type_id
JOIN borehole_organization AS bo
	ON o.organization_id = bo.organization_id
WHERE bo.borehole_id = $1
ORDER BY o.name
