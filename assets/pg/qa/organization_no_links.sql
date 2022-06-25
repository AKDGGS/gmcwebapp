SELECT organization_id AS "Organization ID",
	name AS "Organization Name"
FROM organization
WHERE organization_id NOT IN (
	SELECT DISTINCT organization_id
	FROM borehole_organization
) AND organization_id NOT IN (
	SELECT DISTINCT organization_id
	FROM outcrop_organization
) AND organization_id NOT IN (
	SELECT DISTINCT organization_id
	FROM person_organization
) AND organization_id NOT IN (
	SELECT DISTINCT organization_id
	FROM publication_organization
) AND organization_id NOT IN (
	SELECT DISTINCT organization_id
	FROM well_stratigraphy_organization
) AND organization_id NOT IN (
	SELECT DISTINCT organization_id
	FROM well_operator
) AND organization_id NOT IN (
	SELECT DISTINCT organization_id
	FROM project
) AND organization_id NOT IN (
	SELECT DISTINCT organization_id
	FROM collection
)
