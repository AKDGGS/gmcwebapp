SELECT outcrop_id AS "Outcrop ID", name AS "Outcrop Name"
FROM outcrop
WHERE outcrop_id NOT IN (
	SELECT DISTINCT io.outcrop_id
	FROM inventory_outcrop AS io
	JOIN inventory AS i
		ON i.active AND i.inventory_id = io.inventory_id
)
ORDER BY LOWER(name)
