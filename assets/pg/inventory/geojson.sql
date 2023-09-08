SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', q.features
) AS geojson
FROM ((
	-- Borehole Point
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb
	)) AS features
	FROM inventory AS i
	JOIN inventory_borehole AS ib
		ON ib.inventory_id = i.inventory_id
	JOIN borehole_point AS bp
		ON bp.borehole_id = ib.borehole_id
	JOIN point AS p
		ON p.point_id = bp.point_id
	WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

	) UNION ALL (

	-- Outcrop Point
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb
	)) AS features
	FROM inventory AS i
	JOIN inventory_outcrop AS io
		ON io.inventory_id = i.inventory_id
	JOIN outcrop_point AS op
		ON op.outcrop_id = io.outcrop_id
	JOIN point AS p
		ON p.point_id = op.point_id
	WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

	) UNION ALL (

	-- Outcrop PLSS
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb
	)) AS features
	FROM inventory AS i
	JOIN inventory_outcrop AS io
		ON io.inventory_id = i.inventory_id
	JOIN outcrop_plss AS op
		ON op.outcrop_id = io.outcrop_id
	JOIN plss AS p
		ON p.plss_id = op.plss_id
	WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

	) UNION ALL (

	-- Outcrop place
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb
	)) AS features
	FROM inventory AS i
	JOIN inventory_outcrop AS io
		ON io.inventory_id = i.inventory_id
	JOIN outcrop_place AS op
		ON op.outcrop_id = io.outcrop_id
	JOIN place AS p
		ON p.place_id = op.place_id
	WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

	) UNION ALL (

	-- Outcrop Region
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb
	)) AS features
	FROM inventory AS i
	JOIN inventory_outcrop AS io
		ON io.inventory_id = i.inventory_id
	JOIN outcrop_region AS op
		ON op.outcrop_id = io.outcrop_id
	JOIN region AS p
		ON p.region_id = op.region_id
	WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

	) UNION ALL (

	-- Outcrop Quadrangle
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(q.geog, 5, 0)::jsonb
	)) AS features
	FROM inventory AS i
	JOIN inventory_outcrop AS io
		ON io.inventory_id = i.inventory_id
	JOIN outcrop_quadrangle AS oq
		ON oq.outcrop_id = io.outcrop_id
	JOIN quadrangle AS q
		ON q.quadrangle_id = oq.quadrangle_id
	WHERE i.inventory_id = $1 AND q.geog IS NOT NULL

	) UNION ALL (

	-- Shotline
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb
	)) AS features
	FROM inventory AS i
	JOIN inventory_shotpoint AS isp
		ON isp.inventory_id = i.inventory_id
	JOIN shotpoint AS sp
		ON sp.shotpoint_id = isp.shotpoint_id
	JOIN shotline AS sl
		ON sl.shotline_id = sp.shotline_id
	JOIN shotpoint_point AS spp
		ON spp.shotpoint_id = sp.shotpoint_id
	JOIN point AS p
		ON p.point_id = spp.point_id
	WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

	) UNION ALL (

	-- Well Point
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb
	)) AS features
	FROM inventory AS i
	JOIN inventory_well AS iw
		ON iw.inventory_id = i.inventory_id
	JOIN well_point AS wp
		ON wp.well_id = iw.well_id
	JOIN point AS p
		ON p.point_id = wp.point_id
	WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

	) UNION ALL (

	-- Well Place
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(pl.geog, 5, 0)::jsonb
	)) AS features
	FROM inventory AS i
	JOIN inventory_well AS iw
		ON iw.inventory_id = i.inventory_id
	JOIN well_place AS wpl
		ON wpl.well_id = iw.well_id
	JOIN place AS pl
		ON pl.place_id = wpl.place_id
	WHERE i.inventory_id = $1 AND pl.geog IS NOT NULL

	) UNION ALL (

	-- Well Region
	SELECT jsonb_agg(jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(r.geog, 5, 0)::jsonb
	)) AS features
	FROM inventory AS i
	JOIN inventory_well AS iw
		ON iw.inventory_id = i.inventory_id
	JOIN well_region AS wr
		ON wr.well_id = iw.well_id
	JOIN region AS r
		ON r.region_id = wr.region_id
	WHERE i.inventory_id = $1 AND r.geog IS NOT NULL
)) AS q
WHERE q.features IS NOT NULL
