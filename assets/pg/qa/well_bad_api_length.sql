SELECT well_id AS "Well ID", name AS "Well Name"
FROM well
WHERE api_number IS NOT NULL
	AND LENGTH(api_number) <> 14
