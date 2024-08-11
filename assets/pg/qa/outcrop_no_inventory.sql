WITH ids AS (
	SELECT DISTINCT io.outcrop_id
	FROM inventory_outcrop AS io
	JOIN inventory AS iv
		ON iv.inventory_id = io.inventory_id
	WHERE iv.active
)
SELECT o.outcrop_id AS "Outcrop ID",
	o.name AS "Outcrop Name"
FROM outcrop AS o
LEFT JOIN ids ON ids.outcrop_id = o.outcrop_id
WHERE ids.outcrop_id IS NULL
