SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', q.features
) AS geojson
FROM (
	SELECT jsonb_agg(jsonb_build_object(
			'type', 'Feature',
			'id', b.borehole_id,
			'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb
		) ORDER BY LOWER(b.name)
	) AS features
	FROM borehole AS b
	JOIN borehole_point AS bp
		ON bp.borehole_id = b.borehole_id
	JOIN point AS p
		ON p.point_id = bp.point_id
	WHERE b.borehole_id = $1
) AS q
WHERE q.features IS NOT NULL
