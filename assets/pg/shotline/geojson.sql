SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', jsonb_build_array(
		jsonb_build_object(
			'type', 'Feature',
			'geometry',
			ST_AsGeoJSON(
				CASE
					WHEN ST_NumPoints(q.geog) > 1 THEN ST_Simplify(q.geog,0.0001)
					ELSE ST_PointN(q.geog, 1)
				END, 5, 0)::jsonb,
			'properties', jsonb_strip_nulls(jsonb_build_object(
				'shotline_id', q.shotline_id,
				'name', q.name,
				'alt_names', q.alt_names,
				'year', q.year,
				'remark', q.remark,
				'shotpoint_min', q.shotpoint_min,
				'shotpoint_max', q.shotpoint_max,
			'nearby_shotlines', (
				SELECT jsonb_agg(jsonb_build_object(
					'shotline_id', nearby_shotlines.shotline_id,
					'name', nearby_shotlines.name,
					'distance', nearby_shotlines.distance
				))
				FROM (
					SELECT sl.shotline_id, sl.name,
					ROUND((ST_Distance(q.geog, sl.geog)/ 1609.344):: numeric, 2) AS distance
					FROM (
						SELECT ST_Makeline(
							p.geog::geometry ORDER BY sp.shotpoint_number DESC
						) AS geog, sl.shotline_id, sl.name
						FROM shotline AS sl
						JOIN shotpoint AS sp
							ON sp.shotline_id = sl.shotline_id
						JOIN shotpoint_point AS spp
							ON spp.shotpoint_id = sp.shotpoint_id
						JOIN point AS p
							ON p.point_id = spp.point_id
						WHERE ST_DWithin(q.geog, p.geog, 2414.016)
						GROUP BY sl.shotline_id
					) AS sl
					WHERE q.shotline_id != sl.shotline_id AND ST_DWithin(q.geog, sl.geog, 2414.016)
					ORDER BY distance
					LIMIT 10
				) AS nearby_shotlines
				)
			))
		)
	)) AS geojson
FROM (
	SELECT ST_Makeline(
		p.geog::geometry ORDER BY sp.shotpoint_number DESC
	) AS geog,
		sl.shotline_id AS shotline_id,
		sl.name AS name,
		sl.alt_names as alt_names,
		sl.year as year,
		sl.remark as remark,
		MIN(sp.shotpoint_number) AS shotpoint_min,
		MAX(sp.shotpoint_number) AS shotpoint_max
	FROM shotline AS sl
	JOIN shotpoint AS sp
		ON sp.shotline_id = sl.shotline_id
	JOIN shotpoint_point AS spp
		ON spp.shotpoint_id = sp.shotpoint_id
	JOIN point AS p
		ON p.point_id = spp.point_id
	WHERE sl.shotline_id = $1
	GROUP BY sl.shotline_id, sl.name
) AS q
WHERE q.geog IS NOT NULL
