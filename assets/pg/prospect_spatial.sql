SELECT ST_AsGeoJSON(
	ST_Collect(geog::GEOMETRY),
	5, 0
) AS geojson
FROM borehole AS b
JOIN borehole_point AS bp
	ON bp.borehole_id = b.borehole_id
JOIN point AS p
	ON p.point_id = bp.point_id
WHERE b.prospect_id = $1
