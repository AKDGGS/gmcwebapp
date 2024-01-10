SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', jsonb_agg( jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(o1.geog, 5, 0):: jsonb,
		'properties', jsonb_strip_nulls(jsonb_build_object(
			'outcrop_id', o1.outcrop_id,
			'name', o1.name,
			'nearby_outcrops', (
				SELECT jsonb_agg(jsonb_build_object(
					'outcrop_id', nearby_outcrops.outcrop_id,
					'name', nearby_outcrops.name, 'distance',
					nearby_outcrops.distance
				))
				FROM (
					SELECT
					o2.outcrop_id, o2.name,
					ROUND((ST_Distance(o1.geog, o2.geog)/ 1609.344):: numeric, 2) AS distance
					FROM (
						SELECT o1.outcrop_id, o1.name, p1.geog
						FROM outcrop AS o1
						JOIN outcrop_point AS op1 ON op1.outcrop_id = o1.outcrop_id
						JOIN point AS p1 ON p1.point_id = op1.point_id

						UNION ALL

						SELECT o1.outcrop_id, o1.name, place.geog
						FROM outcrop AS o1
						JOIN outcrop_place AS oplace ON oplace.outcrop_id = o1.outcrop_id
						JOIN place ON place.place_id = oplace.place_id
					) AS o2
					WHERE o1.outcrop_id != o2.outcrop_id AND ST_DWithin(o1.geog, o2.geog, 2414.016)
					ORDER BY distance
					LIMIT 10)
					AS nearby_outcrops
				))
			))
		)
	) AS geojson
FROM (
	SELECT o1.outcrop_id, o1.name, p1.geog
	FROM outcrop AS o1
	JOIN outcrop_point AS op1 ON op1.outcrop_id = o1.outcrop_id
	JOIN point AS p1 ON p1.point_id = op1.point_id

	UNION ALL

	SELECT
	o1.outcrop_id, o1.name, plss.geog
	FROM outcrop AS o1
	JOIN outcrop_plss AS oplss ON oplss.outcrop_id = o1.outcrop_id
	JOIN plss ON plss.plss_id = oplss.plss_id

	UNION ALL

	SELECT o1.outcrop_id, o1.name, place.geog
	FROM outcrop AS o1
	JOIN outcrop_place AS oplace ON oplace.outcrop_id = o1.outcrop_id
	JOIN place ON place.place_id = oplace.place_id

	UNION ALL

	SELECT o1.outcrop_id, o1.name, region.geog
	FROM outcrop AS o1
	JOIN outcrop_region AS oregion ON oregion.outcrop_id = o1.outcrop_id
	JOIN region ON region.region_id = oregion.region_id

	UNION ALL

	SELECT o1.outcrop_id, o1.name, quadrangle.geog
	FROM outcrop AS o1
	JOIN outcrop_quadrangle AS oquad ON oquad.outcrop_id = o1.outcrop_id
	JOIN quadrangle ON quadrangle.quadrangle_id = oquad.quadrangle_id
) AS o1
WHERE o1.outcrop_id = $1 AND o1.geog IS NOT NULL;
