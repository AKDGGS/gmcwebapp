SELECT borehole_id AS "Borehole",
	name AS "Borehole Name"
FROM borehole
WHERE borehole_id NOT IN (
	SELECT borehole_id
	FROM borehole_point
)
ORDER BY borehole_id ASC
