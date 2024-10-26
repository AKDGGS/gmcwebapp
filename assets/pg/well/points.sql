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
			'number', x.well_number
		))
	)) AS features
	FROM (
		SELECT DISTINCT ON (well_id) *
		FROM (
			SELECT w.well_id, w.name, w.well_number, p.geog
			FROM well AS w
			JOIN well_point AS wp ON wp.well_id = w.well_id
			JOIN point AS p ON p.point_id = wp.point_id

			UNION ALL

			SELECT w.well_id, w.name, w.well_number, p.geog
			FROM well AS w
			JOIN well_place AS wp ON wp.well_id = w.well_id
			JOIN place AS p ON p.place_id = wp.place_id
		) AS z
	) AS x
	WHERE x.geog IS NOT NULL
) AS y
WHERE features IS NOT NULL
