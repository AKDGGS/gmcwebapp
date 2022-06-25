SELECT outcrop_id AS "Outcrop ID",
	name AS "Outcrop Name"
FROM outcrop AS o
WHERE outcrop_id NOT IN (
	SELECT DISTINCT io.outcrop_id
	FROM inventory_outcrop AS io
	JOIN inventory AS iv
		ON iv.inventory_id = io.inventory_id
	WHERE iv.active
)
ORDER BY outcrop_id ASC
