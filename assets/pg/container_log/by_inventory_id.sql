SELECT inventory_container_log_id AS id, destination, log_date AS date
FROM inventory_container_log
WHERE inventory_id = $1
ORDER BY log_date DESC
