SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', q.features
) AS geojson
FROM (
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(p1.geog, 5, 0)::jsonb,
		'properties', jsonb_strip_nulls(jsonb_build_object(
			'borehole_id', b1.borehole_id,
			'name', b1.name,
			'alt_names', b1.alt_names,
			'completion_date', to_char(b1.completion_date, 'MM/DD/YY'),
			'onshore', b1.is_onshore,
			'nearby_boreholes', (
				SELECT jsonb_agg(nearby_boreholes)
				FROM (
					SELECT jsonb_build_object(
						'borehole_id', b2.borehole_id,
						'name', b2.name,
						'distance', ROUND((ST_Distance(p1.geog, p2.geog)/1609.344)::numeric, 2)  -- meters to mile conversion
					) AS nearby_boreholes
						FROM borehole AS b2
						JOIN borehole_point AS bp2
							ON bp2.borehole_id = b2.borehole_id
						JOIN point AS p2
							ON p2.point_id = bp2.point_id
						WHERE ST_DWithin(p1.geog, p2.geog, 2414.016) AND b2.borehole_id != b1.borehole_id -- 1.5 (distance threshold) * 1609.344
						ORDER BY ST_Distance(p1.geog, p2.geog)
						LIMIT 10
				) sub
			)
		))
	) ORDER BY LOWER(b1.name)
	) AS features
	FROM borehole AS b1
	JOIN borehole_point AS bp1
		ON bp1.borehole_id = b1.borehole_id
	JOIN point AS p1
		ON p1.point_id = bp1.point_id
	WHERE b1.borehole_id = $1
) AS q
WHERE q.features IS NOT NULL
