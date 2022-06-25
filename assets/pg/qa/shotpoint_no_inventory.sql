SELECT sp.shotpoint_id AS "Shotpoint ID",
	sl.name AS "Shotline",
	sp.shotpoint_number::TEXT AS "Shotpoint Number"
FROM shotpoint AS sp
JOIN shotline AS sl
	ON sl.shotline_id = sp.shotline_id
WHERE shotpoint_id NOT IN (
	SELECT DISTINCT iss.shotpoint_id
	FROM inventory_shotpoint AS iss
	JOIN inventory AS i
		ON i.inventory_id = iss.inventory_id
	WHERE i.active
)
