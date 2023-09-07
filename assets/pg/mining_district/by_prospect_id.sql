SELECT DISTINCT md.mining_district_id AS id,
	md.name
FROM borehole AS b
JOIN borehole_point AS bp
	ON bp.borehole_id = b.borehole_id
JOIN point AS p
	ON p.point_id = bp.point_id
JOIN mining_district AS md
	ON ST_Intersects(md.geog, p.geog)
WHERE b.prospect_id = $1
