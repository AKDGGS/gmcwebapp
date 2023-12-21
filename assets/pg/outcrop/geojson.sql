SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', q.features
) AS geojson
FROM ((
	-- Outcrop Point
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb,
		'properties', jsonb_strip_nulls(jsonb_build_object(
			'outcrop_id', o.outcrop_id,
			'name', o.name,
			'number', o.outcrop_number,
			'year', o.year,
			'onshore', o.is_onshore
		))
	)) AS features
	FROM outcrop AS o
	JOIN outcrop_point AS op ON op.outcrop_id = o.outcrop_id
	JOIN point AS p ON p.point_id = op.point_id
	WHERE o.outcrop_id = $1 AND p.geog IS NOT NULL

	) UNION ALL (

	-- Outcrop PLSS
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb,
		'properties', jsonb_strip_nulls(jsonb_build_object(
			'outcrop_id', o.outcrop_id,
			'name', o.name,
			'number', o.outcrop_number,
			'year', o.year,
			'onshore', o.is_onshore
		))
	)) AS features
	FROM outcrop AS o
	JOIN outcrop_plss AS op ON op.outcrop_id = o.outcrop_id
	JOIN plss AS p ON p.plss_id = op.plss_id
	WHERE o.outcrop_id = $1 AND p.geog IS NOT NULL

	) UNION ALL (

	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb,
		'properties', jsonb_strip_nulls(jsonb_build_object(
			'outcrop_id', o.outcrop_id,
			'name', o.name,
			'number', o.outcrop_number,
			'year', o.year,
			'onshore', o.is_onshore
		))
	)) AS features
	FROM outcrop AS o
	JOIN outcrop_place AS op ON op.outcrop_id = o.outcrop_id
	JOIN place AS p ON p.place_id = op.place_id
	WHERE o.outcrop_id = $1 AND p.geog IS NOT NULL

	) UNION ALL (

	-- Outcrop Region
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb,
		'properties', jsonb_strip_nulls(jsonb_build_object(
			'outcrop_id', o.outcrop_id,
			'name', o.name,
			'number', o.outcrop_number,
			'year', o.year,
			'onshore', o.is_onshore
		))
	)) AS features
	FROM outcrop AS o
	JOIN outcrop_region AS op ON op.outcrop_id = o.outcrop_id
	JOIN region AS p ON p.region_id = op.region_id
	WHERE o.outcrop_id = $1 AND p.geog IS NOT NULL

	) UNION ALL (

	-- Outcrop Quadrangle
	SELECT jsonb_agg(jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(q.geog, 5, 0)::jsonb,
			'properties', jsonb_strip_nulls(jsonb_build_object(
				'outcrop_id', o.outcrop_id,
				'name', o.name,
				'number', o.outcrop_number,
				'year', o.year,
				'onshore', o.is_onshore
			))
	)) AS features
	FROM outcrop AS o
	JOIN outcrop_quadrangle AS oq ON oq.outcrop_id = o.outcrop_id
	JOIN quadrangle AS q ON q.quadrangle_id = oq.quadrangle_id
	WHERE o.outcrop_id = $1 AND q.geog IS NOT NULL
	)
) AS q
WHERE q.features IS NOT NULL
