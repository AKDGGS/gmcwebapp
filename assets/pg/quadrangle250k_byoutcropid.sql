SELECT DISTINCT q.quadrangle_id, q.name
FROM outcrop AS o
JOIN outcrop_point AS op
	ON op.outcrop_id = o.outcrop_id
JOIN point AS p
	ON p.point_id = op.point_id
JOIN quadrangle AS q
	ON ST_Intersects(q.geog, p.geog)
WHERE o.outcrop_id = $1 AND q.scale = 250000
