SELECT sl.shotline_id AS id,
	sl.name,
	sl.alt_names,
	sl.year,
	sl.remark,
FROM inventory_shotpoint AS isp
JOIN shotpoint AS sp
	ON sp.shotpoint_id = isp.shotpoint_id
LEFT OUTER JOIN shotline AS sl
	ON sl.shotline_id = sp.shotline_id
WHERE isp.inventory_id = $1
ORDER BY sl.shotline_id, sp.shotpoint_number
