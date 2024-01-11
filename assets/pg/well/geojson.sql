SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', jsonb_agg( jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(w1.geog, 5, 0):: jsonb,
		'properties', jsonb_strip_nulls(jsonb_build_object(
			'well_id', w1.well_id,
			'name', w1.name,
			'nearby_wells', (
				SELECT jsonb_agg(jsonb_build_object(
					'well_id', nearby_wells.well_id,
					'name', nearby_wells.name, 'distance',
					nearby_wells.distance
				))
				FROM (
					SELECT
					w2.well_id, w2.name,
					ROUND((ST_Distance(w1.geog, w2.geog)/ 1609.344):: numeric, 2) AS distance
					FROM (
						SELECT w1.well_id, w1.name, p1.geog
						FROM well AS w1
						JOIN well_point AS wp1 ON wp1.well_id = w1.well_id
						JOIN point AS p1 ON p1.point_id = wp1.point_id

						UNION ALL

						SELECT w1.well_id, w1.name, place.geog
						FROM well AS w1
						JOIN well_place AS wplace ON wplace.well_id = w1.well_id
						JOIN place ON place.place_id = wplace.place_id
					) AS w2
					WHERE w1.well_id != w2.well_id AND ST_DWithin(w1.geog, w2.geog, 2414.016)
					ORDER BY distance
					LIMIT 10)
					AS nearby_wells
				))
			))
		)
	) AS geojson
FROM (
	SELECT w1.well_id, w1.name, p1.geog
	FROM well AS w1
	JOIN well_point AS wp1 ON wp1.well_id = w1.well_id
	JOIN point AS p1 ON p1.point_id = wp1.point_id

	UNION ALL

	SELECT w1.well_id, w1.name, place.geog
	FROM well AS w1
	JOIN well_place AS wplace ON wplace.well_id = w1.well_id
	JOIN place ON place.place_id = wplace.place_id

	UNION ALL

	SELECT w1.well_id, w1.name, region.geog
	FROM well AS w1
	JOIN well_region AS wregion ON wregion.well_id = w1.well_id
	JOIN region ON region.region_id = wregion.region_id
) AS w1
WHERE w1.well_id = $1 AND w1.geog IS NOT NULL;
