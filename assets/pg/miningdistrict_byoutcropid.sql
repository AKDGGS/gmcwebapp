SELECT DISTINCT md.mining_district_id, md.name
FROM outcrop AS o
JOIN outcrop_point AS op
	ON op.outcrop_id = o.outcrop_id
JOIN point AS p
	ON p.point_id = op.point_id
JOIN mining_district AS md
	ON ST_Intersects(md.geog, p.geog)
WHERE o.outcrop_id = $1
