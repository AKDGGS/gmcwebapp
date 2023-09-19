SELECT 	sp.shotpoint_id AS id,
	sp.shotpoint_number AS number,
	sl.shotline_id AS "shotline.id",
	sl.name AS "shotline.name",
	sl.alt_names AS "shotline.alt_names",
	sl.year AS "shotline.year",
	sl.remark AS "shotline.remark"
FROM inventory_shotpoint AS isp
JOIN shotpoint AS sp
	ON sp.shotpoint_id = isp.shotpoint_id
LEFT OUTER JOIN shotline AS sl
	ON sl.shotline_id = sp.shotline_id
WHERE isp.inventory_id = $1
ORDER BY sl.shotline_id, sp.shotpoint_number