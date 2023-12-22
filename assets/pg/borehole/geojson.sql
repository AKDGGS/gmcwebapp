SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', q.features
) AS geojson
FROM (
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'id', b.borehole_id,
		'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb,
		'properties', jsonb_strip_nulls(jsonb_build_object(
			'borehole_id', b.borehole_id,
			'name', b.name,
			'alt_names', b.alt_names,
			'completion_date', to_char(b.completion_date, 'MM/DD/YY'),
			'onshore', b.is_onshore,
			'nearby_boreholes', (
				SELECT jsonb_agg(nearby_boreholes)
				FROM (
					SELECT jsonb_build_object(
						'borehole_id', b2.borehole_id,
						'name', b2.name,
						'distance', ROUND((ST_Distance(p.geog, p2.geog)/1609.344)::numeric, 2)
					) AS nearby_boreholes
					FROM borehole AS b2
					JOIN borehole_point AS bp2
						ON bp2.borehole_id = b2.borehole_id
					JOIN point AS p2
						ON p2.point_id = bp2.point_id
					WHERE ST_DWithin(p.geog, p2.geog, 1.5 * 1609.344) AND b2.borehole_id != b.borehole_id
					ORDER BY ST_Distance(p.geog, p2.geog)
					LIMIT 10
				) sub
			)
		))
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
