SELECT DISTINCT q.quadrangle_id, q.name
FROM borehole AS b
JOIN borehole_point AS bp
	ON bp.borehole_id = b.borehole_id
JOIN point AS p
	ON p.point_id = bp.point_id
JOIN quadrangle AS q
	ON ST_Intersects(q.geog, p.geog)
WHERE b.prospect_id = $1 AND q.scale = 250000
