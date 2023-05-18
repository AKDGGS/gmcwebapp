SELECT sl.shotline_id AS id,
	sl.name AS name,
	sl.alt_names AS altNames,
	sl.year AS year,
	sl.remark AS remark,
	MIN(sp.shotpoint_number) AS shotpointMin,
	MAX(sp.shotpoint_number) AS shotpointMax
FROM shotline AS sl
JOIN shotpoint AS sp
	ON sp.shotline_id = sl.shotline_id
WHERE sl.shotline_id = $1
GROUP BY sl.shotline_id
