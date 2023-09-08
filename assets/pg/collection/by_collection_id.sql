SELECT collection_id AS id,
	name,
	description,
	organization_id AS "organization.id"
FROM collection c
WHERE collection_id = $1
