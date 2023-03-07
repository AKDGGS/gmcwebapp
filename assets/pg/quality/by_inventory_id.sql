SELECT inventory_quality_id AS "ID",
	check_date as "Date",
	remark,
	username,
	issues::text[]
FROM inventory_quality
WHERE inventory_id = $1
ORDER BY check_date DESC
