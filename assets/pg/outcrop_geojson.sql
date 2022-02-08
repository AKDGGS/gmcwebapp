SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', q.features
) AS geojson
FROM (
	SELECT jsonb_agg(jsonb_build_object(
			'type', 'Feature',
			'id', o.outcrop_id,
			'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb
		) ORDER BY LOWER(o.name)
	) AS features
	FROM outcrop AS o
	JOIN outcrop_point AS op
		ON op.outcrop_id = o.outcrop_id
	JOIN point AS p
		ON p.point_id = op.point_id
	WHERE o.outcrop_id = $1
) AS q
WHERE q.features IS NOT NULL
