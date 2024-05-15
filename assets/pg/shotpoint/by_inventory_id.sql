SELECT
	sp.shotpoint_id AS id,
	sp.shotpoint_number AS "number",
	jsonb_build_object(
		'id', sl.shotline_id,
		'name', sl.name,
		'alt_names', sl.alt_names,
		'year', sl.year,
		'remark', sl.remark
	) AS shotline
FROM inventory_shotpoint AS isp
JOIN shotpoint AS sp
	ON sp.shotpoint_id = isp.shotpoint_id
LEFT OUTER JOIN shotline AS sl
	ON sl.shotline_id = sp.shotline_id
WHERE isp.inventory_id = $1
ORDER BY sl.shotline_id, sp.shotpoint_number
