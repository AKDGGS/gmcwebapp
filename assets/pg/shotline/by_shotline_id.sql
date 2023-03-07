SELECT sl.shotline_id AS "ID",
	sl.name AS "Name",
	sl.alt_names AS "AltNames",
	sl.year AS "Year",
	sl.remark AS "Remark",
	MIN(sp.shotpoint_number) AS "ShotpointMin",
	MAX(sp.shotpoint_number) AS "ShotpointMax"
FROM shotline AS sl
JOIN shotpoint AS sp
	ON sp.shotline_id = sl.shotline_id
WHERE sl.shotline_id = $1
GROUP BY sl.shotline_id
