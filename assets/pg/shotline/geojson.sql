SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', jsonb_build_array(jsonb_build_object(
		'type', 'Feature', 'geometry',
		ST_AsGeoJSON(
			CASE
				WHEN ST_NumPoints(q.geog) > 1 THEN ST_Simplify(q.geog,0.0001)
				ELSE ST_PointN(q.geog, 1)
			END
		, 5, 0)::jsonb
	))
) AS geojson
FROM (
	SELECT ST_Makeline(
		p.geog::geometry ORDER BY sp.shotpoint_number DESC
	) AS geog
	FROM shotline AS s
	JOIN shotpoint AS sp
		ON sp.shotline_id = s.shotline_id
	JOIN shotpoint_point AS spp
		ON spp.shotpoint_id = sp.shotpoint_id
	JOIN point AS p
		ON p.point_id = spp.point_id
	WHERE s.shotline_id = $1
) AS q
WHERE q.geog IS NOT NULL
