SELECT
	inventory_container_log_id AS id,
	log_date AS "date",
	destination
FROM inventory_container_log
WHERE inventory_id = $1
ORDER BY log_date DESC
