SELECT well_id AS "Well ID", name AS "Well Name"
FROM well
WHERE (
	well_number IS NULL OR LENGTH(TRIM(BOTH FROM well_number)) = 0
) AND well_id IN (
	SELECT DISTINCT well_id
	FROM inventory_well AS iw
	JOIN inventory AS i
		ON i.inventory_id = iw.inventory_id
	WHERE active
)
