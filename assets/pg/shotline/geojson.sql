SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', jsonb_build_array(jsonb_build_object(
		'type', 'Feature', 'geometry',
		ST_AsGeoJSON(q.geog, 5, 0)::jsonb
	))
) AS geojson
FROM (
	SELECT ST_Makeline(
		p.geog::geometry ORDER BY sp.shotpoint_number DESC
	) AS geog
	FROM shotline AS s
	JOIN shotpoint AS sp
		ON sp.shotline_id = s.shotline_id
	JOIN shotpoint_point as spp
		ON spp.shotpoint_id = sp.shotpoint_id
	JOIN point AS p
		ON p.point_id = spp.point_id
	WHERE s.shotline_id = $1
) AS q
WHERE q.geog IS NOT NULL
