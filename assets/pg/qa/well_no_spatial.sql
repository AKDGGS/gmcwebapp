SELECT well_id AS "Well ID", name AS "Well Name"
FROM well
WHERE well_id NOT IN (
	SELECT well_id
	FROM well_point
) AND well_id NOT IN (
	SELECT well_id
	FROM well_place AS wp
	JOIN place AS pl
		ON pl.place_id = wp.place_id
	WHERE pl.geog IS NOT NULL
) AND well_id NOT IN (
	SELECT well_id
	FROM well_region AS wr
	JOIN region AS r
		ON wr.region_id = r.region_id
	WHERE r.geog IS NOT NULL
)
ORDER BY well_id
