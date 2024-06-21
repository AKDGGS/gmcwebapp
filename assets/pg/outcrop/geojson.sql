SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', q.features
) AS geojson
FROM (
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(geometry, 5, 0)::jsonb,
		'properties', jsonb_strip_nulls(jsonb_build_object(
			'outcrop_id', outcrop_id,
			'name', name,
			'nearby_outcrops', (
				SELECT jsonb_agg(jsonb_build_object(
					'outcrop_id', nearby_outcrops.outcrop_id,
					'name', nearby_outcrops.name,
					'distance', nearby_outcrops.distance
				))
				FROM (
					SELECT
						o1_copy.outcrop_id, o1_copy.name,
						ROUND((ST_Distance(geometry, o1_copy.geog) / 1609.344)::numeric, 2) AS distance
					FROM (
						SELECT o1_copy.outcrop_id, o1_copy.name, p1.geog
						FROM outcrop AS o1_copy
						JOIN outcrop_point AS op1 ON op1.outcrop_id = o1_copy.outcrop_id
						JOIN point AS p1 ON p1.point_id = op1.point_id
						WHERE combined_data.outcrop_id != o1_copy.outcrop_id AND ST_DWithin(geometry, p1.geog, 2414.016)
					) AS o1_copy

					UNION ALL

					SELECT
						o2.outcrop_id, o2.name,
						ROUND((ST_Distance(geometry, o2.geog) / 1609.344)::numeric, 2) AS distance
					FROM (
						SELECT o2.outcrop_id, o2.name, place.geog
						FROM outcrop AS o2
						JOIN outcrop_place AS oplace ON oplace.outcrop_id = o2.outcrop_id
						JOIN place ON place.place_id = oplace.place_id
						WHERE combined_data.outcrop_id != o2.outcrop_id AND ST_DWithin(geometry, place.geog, 2414.016)
					) AS o2
					ORDER BY distance
					LIMIT 10
				) AS nearby_outcrops
			)
		))
	)) AS features
	FROM (
		SELECT p1.geog AS geometry, o1.outcrop_id, o1.name
		FROM outcrop AS o1
		JOIN outcrop_point AS op1 ON op1.outcrop_id = o1.outcrop_id
		JOIN point AS p1 ON p1.point_id = op1.point_id

		UNION ALL

		SELECT plss.geog AS geometry, o1.outcrop_id, o1.name
		FROM outcrop AS o1
		JOIN outcrop_plss AS oplss ON oplss.outcrop_id = o1.outcrop_id
		JOIN plss ON plss.plss_id = oplss.plss_id

		UNION ALL

		SELECT place.geog AS geomerty, o1.outcrop_id, o1.name
		FROM outcrop AS o1
		JOIN outcrop_place AS oplace ON oplace.outcrop_id = o1.outcrop_id
		JOIN place ON place.place_id = oplace.place_id

		UNION ALL

		SELECT region.geog AS geometry, o1.outcrop_id, o1.name
		FROM outcrop AS o1
		JOIN outcrop_region AS oregion ON oregion.outcrop_id = o1.outcrop_id
		JOIN region ON region.region_id = oregion.region_id
		UNION ALL

		SELECT quadrangle.geog AS geometry, o1.outcrop_id, o1.name
		FROM outcrop AS o1
		JOIN outcrop_quadrangle AS oquad ON oquad.outcrop_id = o1.outcrop_id
		JOIN quadrangle ON quadrangle.quadrangle_id = oquad.quadrangle_id
	) AS combined_data
	WHERE combined_data.outcrop_id = $1 AND geometry IS NOT NULL
) AS q
WHERE q.features IS NOT NULL;
