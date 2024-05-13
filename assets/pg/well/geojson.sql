SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', features
)
FROM (
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(x.geog,5,0)::jsonb,
		'properties', jsonb_strip_nulls(jsonb_build_object(
			'well_id', x.well_id,
			'name', x.name,
			'nearby_wells', (
				SELECT jsonb_agg(jsonb_build_object(
					'well_id', nearby_wells.well_id,
					'name', nearby_wells.name,
					'distance', nearby_wells.distance
				)) AS nearby_wells
				FROM (
					SELECT s.well_id, s.name,
					ROUND(
						(ST_Distance(x.geog, s.geog)/1609.344)::numeric, 2
					) AS distance
					FROM (
						SELECT w.well_id, w.name, p.geog
						FROM well AS w
						JOIN well_point AS wp ON wp.well_id = w.well_id
						JOIN point AS p ON p.point_id = wp.point_id

						UNION ALL

						SELECT w.well_id, w.name, p.geog
						FROM well AS w
						JOIN well_place AS wp ON wp.well_id = w.well_id
						JOIN place AS p ON p.place_id = wp.place_id
					) AS s
					WHERE x.well_id != s.well_id
						AND ST_DWithin(x.geog, s.geog, 2414.016)
					ORDER BY distance
					LIMIT 10
				) AS nearby_wells
			)
		))
	)) AS features
	FROM (
		SELECT w.well_id, w.name, p.geog
		FROM well AS w
		JOIN well_point AS wp ON wp.well_id = w.well_id
		JOIN point AS p ON p.point_id = wp.point_id

		UNION ALL

		SELECT w.well_id, w.name, p.geog
		FROM well AS w
		JOIN well_place AS wp ON wp.well_id = w.well_id
		JOIN place AS p ON p.place_id = wp.place_id

		UNION ALL

		SELECT w.well_id, w.name, r.geog
		FROM well AS w
		JOIN well_region AS wr ON wr.well_id = w.well_id
		JOIN region AS r ON r.region_id = wr.region_id
	) AS x
	WHERE x.geog IS NOT NULL AND x.well_id = $1
) AS y
WHERE features IS NOT NULL
