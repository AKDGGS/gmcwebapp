SELECT o.organization_id, o.name, ot.name AS organization_type
FROM organization AS o
JOIN organization_type AS ot
    ON o.organization_type_id = ot.organization_type_id
JOIN outcrop_organization AS oo
    ON o.organization_id = oo.organization_id
WHERE oo.outcrop_id = $1
ORDER BY o.name
