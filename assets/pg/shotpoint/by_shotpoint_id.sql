SELECT 	sp.shotpoint_id AS id,
	sp.shotpoint_number AS number,
	sl.shotline_id AS "shotline.id",
	sl.name AS "shotline.name",
	sl.alt_names AS "shotline.alt_names",
	sl.year AS "shotline.year",
	sl.remark AS "shotline.remark"
FROM shotpoint AS sp
LEFT OUTER JOIN shotline AS sl
	ON sl.shotline_id = sp.shotline_id
WHERE sp.shotpoint_id = ANY($1)
ORDER BY sl.shotline_id, sp.shotpoint_number;
