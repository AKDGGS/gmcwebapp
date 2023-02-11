SELECT o.organization_id AS "ID", o.name, ot.name AS "Type", wo.is_current AS "Current"
FROM organization AS o
JOIN organization_type AS ot
	ON o.organization_type_id = ot.organization_type_id
JOIN well_operator AS wo
	ON o.organization_id = wo.organization_id
WHERE wo.well_id = $1
ORDER BY wo.is_current ASC, o.name
