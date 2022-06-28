SELECT well_id AS "Well ID",
	name AS "Well Name",
	well_number AS "Well Number"
FROM well
WHERE well_id NOT IN (
	SELECT DISTINCT well_id
	FROM inventory_well AS iw
	JOIN inventory AS i
		ON i.inventory_id = iw.inventory_id
	WHERE i.active
)
