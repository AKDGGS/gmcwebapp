-- Materialized spatial view for inventory
CREATE MATERIALIZED VIEW inventory_geog AS (
	-- Boreholes
	SELECT DISTINCT ON (inventory_id, borehole_id)
		inventory_id, geog
	FROM ((
		-- Borehole Point
		SELECT ib.inventory_id, ib.borehole_id, p.geog
		FROM inventory_borehole AS ib
		JOIN borehole_point AS bp ON bp.borehole_id = ib.borehole_id
		JOIN point AS p ON p.point_id = bp.point_id
		WHERE p.geog IS NOT NULL
	)) AS q

	UNION ALL

	-- Outcrops
	SELECT DISTINCT ON (inventory_id, outcrop_id)
		inventory_id, geog
	FROM ((
		-- Outcrop Point
		SELECT io.inventory_id, io.outcrop_id, p.geog
		FROM inventory_outcrop AS io
		JOIN outcrop_point AS op ON op.outcrop_id = io.outcrop_id
		JOIN point AS p ON p.point_id = op.point_id
		WHERE p.geog IS NOT NULL

		) UNION ALL (

		-- Outcrop PLSS
		SELECT io.inventory_id, io.outcrop_id, p.geog
		FROM inventory_outcrop AS io
		JOIN outcrop_plss AS op ON op.outcrop_id = io.outcrop_id
		JOIN plss AS p ON p.plss_id = op.plss_id
		WHERE p.geog IS NOT NULL

		) UNION ALL (

		-- Outcrop Place
		SELECT io.inventory_id, io.outcrop_id, p.geog
		FROM inventory_outcrop AS io
		JOIN outcrop_place AS op ON op.outcrop_id = io.outcrop_id
		JOIN place AS p ON p.place_id = op.place_id
		WHERE p.geog IS NOT NULL

		) UNION ALL (

		-- Outcrop Region
		SELECT io.inventory_id, io.outcrop_id, p.geog
		FROM inventory_outcrop AS io
		JOIN outcrop_region AS op ON op.outcrop_id = io.outcrop_id
		JOIN region AS p ON p.region_id = op.region_id
		WHERE p.geog IS NOT NULL

		) UNION ALL (

		-- Outcrop Quadrangle
		SELECT io.inventory_id, io.outcrop_id, q.geog
		FROM inventory_outcrop AS io
		JOIN outcrop_quadrangle AS oq ON oq.outcrop_id = io.outcrop_id
		JOIN quadrangle AS q ON q.quadrangle_id = oq.quadrangle_id
		WHERE q.geog IS NOT NULL
		ORDER BY q.scale ASC
	)) AS q

	UNION ALL

	-- Wells
	SELECT DISTINCT ON (inventory_id, well_id)
		inventory_id, geog
	FROM ((
		-- Well Point
		SELECT iw.inventory_id, iw.well_id, p.geog
		FROM inventory_well AS iw
		JOIN well_point AS wp ON wp.well_id = iw.well_id
		JOIN point AS p ON p.point_id = wp.point_id
		WHERE p.geog IS NOT NULL

		) UNION ALL (

		-- Well Place
		SELECT iw.inventory_id, iw.well_id, p.geog
		FROM inventory_well AS iw
		JOIN well_place AS wp ON wp.well_id = iw.well_id
		JOIN place AS p ON p.place_id = wp.place_id
		WHERE p.geog IS NOT NULL

		) UNION ALL (

		-- Well Region
		SELECT iw.inventory_id, iw.well_id, p.geog
		FROM inventory_well AS iw
		JOIN well_region AS wp ON wp.well_id = iw.well_id
		JOIN region AS p ON p.region_id = wp.region_id
		WHERE p.geog IS NOT NULL
	)) AS q

	UNION ALL

	-- Shotlines
	SELECT inventory_id,
		CASE WHEN ST_NumPoints(geom) > 1 THEN
			ST_Simplify(geom,0.0001)::GEOGRAPHY
		ELSE
			ST_PointN(geom, 1)::GEOGRAPHY
		END AS geog
	FROM (
		SELECT isp.inventory_id,
			ST_Makeline(geog::geometry ORDER BY sp.shotpoint_number) AS geom
		FROM inventory_shotpoint AS isp
		JOIN shotpoint AS sp ON sp.shotpoint_id = isp.shotpoint_id
		JOIN shotpoint_point AS spp ON spp.shotpoint_id = sp.shotpoint_id
		JOIN point AS p ON p.point_id = spp.point_id
		GROUP BY isp.inventory_id, sp.shotline_id
		ORDER BY isp.inventory_id, sp.shotline_id ASC
	) AS q

	UNION ALL

	-- Publications
	SELECT ip.inventory_id, q.geog
	FROM inventory_publication AS ip
	JOIN publication_quadrangle AS pq
		ON pq.publication_id = ip.publication_id
	JOIN quadrangle AS q
		ON q.quadrangle_id = pq.quadrangle_id
);
