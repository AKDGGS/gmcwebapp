SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', q.features
) AS geojson
FROM ((
	-- Well Point
	SELECT jsonb_agg(jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb
		)
	) AS features
		FROM well AS w
		JOIN well_point AS wp ON wp.well_id = w.well_id
		JOIN point AS p ON p.point_id = wp.point_id
		WHERE w.well_id = $1 AND p.geog IS NOT NULL

	) UNION ALL (

	-- Well Place
	SELECT jsonb_agg(jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(pl.geog, 5, 0)::jsonb
		)
	) AS features
		FROM well AS w
		JOIN well_place AS wpl ON wpl.well_id = w.well_id
		JOIN place AS pl ON pl.place_id = wpl.place_id
		WHERE w.well_id = $1 AND pl.geog IS NOT NULL

	) UNION ALL (

	-- Well Region
	SELECT jsonb_agg(jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(r.geog, 5, 0)::jsonb
		)
	) AS features
		FROM well AS w
		JOIN well_region AS wr ON wr.well_id = w.well_id
		JOIN region AS r ON r.region_id = wr.region_id
		WHERE w.well_id = $1 AND r.geog IS NOT NULL
	)
) AS q
WHERE q.features IS NOT NULL
