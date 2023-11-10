SELECT container_type_id
FROM container_type
WHERE name = $1
LIMIT 1
