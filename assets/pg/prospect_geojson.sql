SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', ((
		SELECT jsonb_agg(jsonb_build_object(
			'type', 'Feature',
			'id', b.borehole_id,
			'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb,
			'properties', jsonb_strip_nulls(jsonb_build_object(
				'borehole_id', b.borehole_id,
				'name', b.name,
				'alt_names', b.alt_names,
				'completion_date', to_char(b.completion_date, 'MM/DD/YY')
			))
		) ORDER BY LOWER(b.name) DESC)
		FROM borehole AS b
		JOIN borehole_point AS bp
			ON bp.borehole_id = b.borehole_id
		JOIN point AS p
			ON p.point_id = bp.point_id
		WHERE b.prospect_id = $1
	))
) as geojson
