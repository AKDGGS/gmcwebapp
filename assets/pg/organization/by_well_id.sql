SELECT o.organization_id AS id, o.name, ot.name AS type, o.remark,
wo.is_current AS current
FROM organization AS o
JOIN organization_type AS ot
	ON o.organization_type_id = ot.organization_type_id
JOIN well_operator AS wo
	ON o.organization_id = wo.organization_id
WHERE wo.well_id = $1
ORDER BY wo.is_current ASC, o.name
