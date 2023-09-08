SELECT sl.shotline_id AS id,
	sl.name,
	sl.alt_names,
	sl.year,
	sl.remark,
	MIN(sp.shotpoint_number) AS shotpointMin,
	MAX(sp.shotpoint_number) AS shotpointMax
FROM shotline AS sl
JOIN shotpoint AS sp
	ON sp.shotline_id = sl.shotline_id
WHERE sl.shotline_id = $1
GROUP BY sl.shotline_id
