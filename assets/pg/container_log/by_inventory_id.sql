SELECT inventory_container_log_id AS "ID", destination, log_date AS "Date"
FROM inventory_container_log
WHERE inventory_id = $1
ORDER BY log_date DESC
