SELECT sp.shotline_id AS id,
	sp.shotpoint_number AS number
FROM shotpoint AS sp
WHERE sp.shotline_id = $1
ORDER BY shotpoint_number;
