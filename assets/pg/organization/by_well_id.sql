SELECT o.organization_id, o.name, ot.name AS operator_type, wo.is_current
FROM organization AS o
JOIN organization_type AS ot
	ON o.organization_type_id = ot.organization_type_id
JOIN well_operator AS wo
	ON o.organization_id = wo.organization_id
WHERE wo.well_id = $1
ORDER BY wo.is_current DESC, o.name
