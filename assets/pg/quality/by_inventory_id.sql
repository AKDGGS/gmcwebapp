SELECT
	inventory_quality_id AS id,
	check_date AS "date",
	remark,
	username,
	issues::text[]
FROM inventory_quality
WHERE inventory_id = $1
ORDER BY check_date DESC
