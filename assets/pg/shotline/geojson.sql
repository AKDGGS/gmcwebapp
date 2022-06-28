SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', q.features
) AS geojson
FROM (
	SELECT jsonb_agg(jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb
		)
	) AS features
	FROM shotline AS s
	JOIN shotpoint AS sp
		ON sp.shotline_id = s.shotline_id
	JOIN shotpoint_point as spp
	ON spp.shotpoint_id = sp.shotpoint_id
	JOIN point AS p
	ON p.point_id = spp.point_id
	WHERE s.shotline_id = $1
) AS q
WHERE q.features IS NOT NULL
