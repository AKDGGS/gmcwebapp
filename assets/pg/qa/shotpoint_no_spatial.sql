SELECT sp.shotpoint_id AS "Shotpoint ID",
	sp.shotpoint_number::TEXT AS "Shotpoint Number",
	sl.name AS "Shotline"	
FROM shotpoint AS sp
LEFT OUTER JOIN shotline AS sl
	ON sp.shotline_id = sl.shotline_id
WHERE sp.shotpoint_id NOT IN (
	SELECT shotpoint_id
	FROM shotpoint_point
)
ORDER BY sl.name
