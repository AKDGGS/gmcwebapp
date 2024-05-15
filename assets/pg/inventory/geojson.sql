SELECT jsonb_build_object(
	'type', 'FeatureCollection',
	'features', x.features
) AS geojson
FROM (
	SELECT jsonb_agg(q.features) AS features
	FROM ((
		-- Borehole Point
		SELECT jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb,
			'properties', jsonb_strip_nulls(jsonb_build_object(
				'borehole_id', b.borehole_id,
				'name', b.name,
				'alt_names', b.alt_names,
				'onshore', b.is_onshore,
				'completion_date', to_char(b.completion_date, 'MM/DD/YY')
			))
		) AS features
		FROM inventory AS i
		JOIN inventory_borehole AS ib
			ON ib.inventory_id = i.inventory_id
		JOIN borehole AS b
			ON ib.borehole_id = b.borehole_id
		JOIN borehole_point AS bp
			ON bp.borehole_id = ib.borehole_id
		JOIN point AS p
			ON p.point_id = bp.point_id
		WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

		) UNION ALL (

		-- Outcrop Point
		SELECT jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb,
			'properties', jsonb_strip_nulls(jsonb_build_object(
				'outcrop_id', o.outcrop_id,
				'name', o.name,
				'number', o.outcrop_number,
				'year', o.year,
				'onshore', o.is_onshore
			))
		) AS features
		FROM inventory AS i
		JOIN inventory_outcrop AS io
			ON io.inventory_id = i.inventory_id
		JOIN outcrop AS o
			ON io.outcrop_id = o.outcrop_id
		JOIN outcrop_point AS op
			ON op.outcrop_id = io.outcrop_id
		JOIN point AS p
			ON p.point_id = op.point_id
		WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

		) UNION ALL (

		-- Outcrop PLSS
		SELECT jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb,
			'properties', jsonb_strip_nulls(jsonb_build_object(
				'outcrop_id', o.outcrop_id,
				'name', o.name,
				'number', o.outcrop_number,
				'year', o.year,
				'onshore', o.is_onshore
			))
		) AS features
		FROM inventory AS i
		JOIN inventory_outcrop AS io
			ON io.inventory_id = i.inventory_id
		JOIN outcrop AS o
			ON io.outcrop_id = o.outcrop_id
		JOIN outcrop_plss AS op
			ON op.outcrop_id = io.outcrop_id
		JOIN plss AS p
			ON p.plss_id = op.plss_id
		WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

		) UNION ALL (

		-- Outcrop place
		SELECT jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb,
			'properties', jsonb_strip_nulls(jsonb_build_object(
				'outcrop_id', o.outcrop_id,
				'name', o.name,
				'number', o.outcrop_number,
				'year', o.year,
				'onshore', o.is_onshore
			))
		) AS features
		FROM inventory AS i
		JOIN inventory_outcrop AS io
			ON io.inventory_id = i.inventory_id
		JOIN outcrop AS o
			ON io.outcrop_id = o.outcrop_id
		JOIN outcrop_place AS op
			ON op.outcrop_id = io.outcrop_id
		JOIN place AS p
			ON p.place_id = op.place_id
		WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

		) UNION ALL (

		-- Outcrop Region
		SELECT jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb,
			'properties', jsonb_strip_nulls(jsonb_build_object(
				'outcrop_id', o.outcrop_id,
				'name', o.name,
				'number', o.outcrop_number,
				'year', o.year,
				'onshore', o.is_onshore
			))
		) AS features
		FROM inventory AS i
		JOIN inventory_outcrop AS io
			ON io.inventory_id = i.inventory_id
		JOIN outcrop AS o
			ON io.outcrop_id = o.outcrop_id
		JOIN outcrop_region AS op
			ON op.outcrop_id = io.outcrop_id
		JOIN region AS p
			ON p.region_id = op.region_id
		WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

		) UNION ALL (

		-- Outcrop Quadrangle
		SELECT jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(q.geog, 5, 0)::jsonb,
			'properties', jsonb_strip_nulls(jsonb_build_object(
				'outcrop_id', o.outcrop_id,
				'name', o.name,
				'number', o.outcrop_number,
				'year', o.year,
				'onshore', o.is_onshore
			))
		) AS features
		FROM inventory AS i
		JOIN inventory_outcrop AS io
			ON io.inventory_id = i.inventory_id
		JOIN outcrop AS o
			ON io.outcrop_id = o.outcrop_id
		JOIN outcrop_quadrangle AS oq
			ON oq.outcrop_id = io.outcrop_id
		JOIN quadrangle AS q
			ON q.quadrangle_id = oq.quadrangle_id
		WHERE i.inventory_id = $1 AND q.geog IS NOT NULL

		) UNION ALL (

		-- Shotline
		SELECT jsonb_build_object(
		'type', 'Feature',
		'geometry', ST_AsGeoJSON(
		CASE
				WHEN ST_NumPoints(line.geog) > 1 THEN ST_Simplify(line.geog,0.0001)
				ELSE ST_PointN(line.geog, 1)
		END
		, 5, 0)::jsonb,
		'properties',
		jsonb_strip_nulls(jsonb_build_object(
			'shotline_id', line.shotline_id,
			'name', line.name,
			'alt_names', line.alt_names,
			'year', line.year,
			'remark', line.remark,
			'shotpoint_min', line.shotpoint_min,
			'shotpoint_max', line.shotpoint_max
		))
	) AS features
	FROM (
		SELECT ST_MakeLine(p.geog::geometry ORDER BY sp.shotpoint_number DESC) AS geog,
			sl.shotline_id AS shotline_id,
			sl.name AS name,
			sl.alt_names as alt_names,
			sl.year as year,
			sl.remark as remark,
			MIN(sp.shotpoint_number) AS shotpoint_min,
			MAX(sp.shotpoint_number) AS shotpoint_max
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
		GROUP BY  sl.shotline_id, sl.name
		) AS line
		JOIN shotline ON line.shotline_id = shotline.shotline_id

		) UNION ALL (

		-- Well Point
		SELECT jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(p.geog, 5, 0)::jsonb,
			'properties', jsonb_strip_nulls(jsonb_build_object(
				'well_id', w.well_id,
				'name', w.name,
				'well_number', w.well_number,
				'api_number', w.api_number,
				'onshore', w.is_onshore,
				'federal', w.is_federal,
				'spud_date', w.spud_date,
				'completion_date', w.completion_date,
				'measured_depth', w.measured_depth,
				'elevation', w.elevation,
				'elevation_kb', w.elevation_kb,
				'permit_status', w.permit_status,
				'permit_number', w.permit_number,
				'completion_status', w.completion_status
			))
		) AS features
		FROM inventory AS i
		JOIN inventory_well AS iw
			ON iw.inventory_id = i.inventory_id
		JOIN well as w
		 ON iw.well_id = w.well_id
		JOIN well_point AS wp
			ON wp.well_id = iw.well_id
		JOIN point AS p
			ON p.point_id = wp.point_id
		WHERE i.inventory_id = $1 AND p.geog IS NOT NULL

		) UNION ALL (

		-- Well Place
		SELECT jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(pl.geog, 5, 0)::jsonb,
			'properties', jsonb_strip_nulls(jsonb_build_object(
				'well_id', w.well_id,
				'name', w.name,
				'well_number', w.well_number,
				'api_number', w.api_number
			))
		) AS features
		FROM inventory AS i
		JOIN inventory_well AS iw
			ON iw.inventory_id = i.inventory_id
		JOIN well as w
		 ON iw.well_id = w.well_id
		JOIN well_place AS wpl
			ON wpl.well_id = iw.well_id
		JOIN place AS pl
			ON pl.place_id = wpl.place_id
		WHERE i.inventory_id = $1 AND pl.geog IS NOT NULL

		) UNION ALL (

		-- Well Region
		SELECT jsonb_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(r.geog, 5, 0)::jsonb,
			'properties', jsonb_strip_nulls(jsonb_build_object(
				'well_id', w.well_id,
				'name', w.name,
				'well_number', w.well_number,
				'api_number', w.api_number
			))
		) AS features
		FROM inventory AS i
		JOIN inventory_well AS iw
			ON iw.inventory_id = i.inventory_id
		JOIN well as w
			ON iw.well_id = w.well_id
		JOIN well_region AS wr
			ON wr.well_id = iw.well_id
		JOIN region AS r
			ON r.region_id = wr.region_id
		WHERE i.inventory_id = $1 AND r.geog IS NOT NULL
	)) AS q
) AS x
WHERE x.features IS NOT NULL
