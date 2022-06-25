SELECT api_number AS "API Number",
	STRING_AGG(well_id::TEXT, ',') AS "Well IDs"
FROM well
WHERE api_number IS NOT NULL
GROUP BY api_number
HAVING COUNT(*) > 1
