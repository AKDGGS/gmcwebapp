SELECT DISTINCT q.quadrangle_id AS id,
	q.name
FROM well AS w
JOIN well_point AS wp
	ON wp.well_id = w.well_id
JOIN point AS p
	ON p.point_id = wp.point_id
JOIN quadrangle AS q
	ON ST_Intersects(q.geog, p.geog)
WHERE w.well_id = $1 AND q.scale = 250000
