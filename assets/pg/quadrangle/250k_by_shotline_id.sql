SELECT DISTINCT q.quadrangle_id AS id,
	q.name
FROM shotline AS s
JOIN shotpoint AS sp
	ON sp.shotline_id = s.shotline_id
JOIN shotpoint_point AS spp
	ON spp.shotpoint_id = sp.shotpoint_id
JOIN point AS p
	ON p.point_id = spp.point_id
JOIN quadrangle AS q
	ON ST_Intersects(q.geog, p.geog)
WHERE s.shotline_id = $1 AND q.scale = 250000
