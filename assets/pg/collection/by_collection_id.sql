SELECT collection_id AS id,
name,
description,
organization_id
FROM collection c
WHERE collection_id = $1
