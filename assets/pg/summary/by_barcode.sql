WITH ivs AS (
	SELECT ARRAY_AGG(inventory_id) AS ids
	FROM inventory
	WHERE active AND (
		barcode = $1
		OR alt_barcode = $1
		OR container_id IN (
			WITH RECURSIVE r AS (
				SELECT container_id
				FROM container
				WHERE barcode = $1
					OR alt_barcode = $1

				UNION ALL

				SELECT co.container_id
				FROM r
				JOIN container AS co
					ON r.container_id = co.parent_container_id
			) SELECT container_id FROM r
		)
	)
), kws AS (
	SELECT
		ARRAY_AGG(DISTINCT kw.keyword ORDER BY kw.keyword) AS keywords
	FROM (
		SELECT UNNEST(i.keywords) AS keyword
		FROM ivs
		JOIN inventory AS i ON i.inventory_id = ANY(ivs.ids)
	) AS kw
), bcs AS (
	SELECT
	ARRAY_AGG(DISTINCT COALESCE(i.barcode, i.alt_barcode)
		ORDER BY COALESCE(i.barcode, i.alt_barcode)) AS barcodes,
	COUNT(DISTINCT COALESCE(i.barcode, i.alt_barcode)) AS barcode_total
	FROM ivs
	JOIN inventory AS i ON i.inventory_id = ANY(ivs.ids)
	WHERE COALESCE(barcode, alt_barcode) IS NOT NULL
), cts AS (
		SELECT
			c.path_cache AS container,
			COUNT(i.inventory_id) AS container_total
		FROM ivs
		JOIN inventory AS i ON i.inventory_id = ANY(ivs.ids)
		JOIN container AS c ON c.container_id = i.container_id
		GROUP BY c.path_cache
	-- ) AS q
), cls AS (
	SELECT
		jsonb_agg(cb) AS collections
	FROM (
		SELECT jsonb_build_object(
			'collection', COALESCE(c.name, 'None'),
			'collection_total', COUNT(i.inventory_id)
		) AS cb
		FROM ivs
		JOIN inventory AS i ON i.inventory_id = ANY(ivs.ids)
		LEFT JOIN collection AS c ON c.collection_id = i.collection_id
		GROUP BY c.name
	) AS q
)
, bhs AS (
	SELECT
		jsonb_agg(bh) AS boreholes
	FROM (
		SELECT jsonb_build_object(
			'prospect', ps.name,
			'borehole', b.name,
			'borehole_total', COUNT(DISTINCT COALESCE(i.barcode, i.alt_barcode))
		) AS bh
		FROM ivs
		JOIN inventory_borehole AS ib ON ib.inventory_id = ANY(ivs.ids)
		JOIN inventory AS i ON i.inventory_id = ib.inventory_id
		LEFT JOIN borehole AS b ON b.borehole_id = ib.borehole_id
		LEFT JOIN prospect AS ps ON ps.prospect_id = b.prospect_id
		GROUP BY b.name, ps.name
		ORDER BY COUNT(ib.inventory_id) DESC, ps.name, b.name
		LIMIT 100
	) AS q
),ocs AS (
	SELECT
		jsonb_agg(oc) AS outcrops
	FROM (
		SELECT jsonb_build_object(
			'outcrop', (o.name || COALESCE(' - ' || o.outcrop_number, '')),
			'outcrop_total', COUNT(DISTINCT COALESCE(i.barcode, i.alt_barcode))
		) AS oc
		FROM ivs
		JOIN inventory_outcrop AS io ON io.inventory_id = ANY(ivs.ids)
		JOIN inventory AS i ON i.inventory_id = io.inventory_id
		LEFT JOIN outcrop AS o ON o.outcrop_id = io.outcrop_id
		GROUP BY o.name, o.outcrop_number
		ORDER BY COUNT(DISTINCT COALESCE(i.barcode, i.alt_barcode))  DESC,
			o.name, o.outcrop_number
			LIMIT 100
	) AS q
)
,sls AS (
	SELECT
		jsonb_agg(sl) AS shotlines
	FROM (
		SELECT jsonb_build_object(
			'shotline', sl.name,
			'shotline_total', COUNT(DISTINCT COALESCE(i.barcode, i.alt_barcode))
		) AS sl
		FROM ivs
		JOIN inventory_shotpoint AS isp ON isp.inventory_id = ANY(ivs.ids)
		JOIN inventory AS i ON i.inventory_id = isp.inventory_id
		JOIN shotpoint AS sp ON sp.shotpoint_id = isp.shotpoint_id
		JOIN shotline sl ON sl.shotline_id = sp.shotline_id
		GROUP BY sl.name
		ORDER BY sl.name
		LIMIT 100
	) AS q
),ws AS (
	SELECT
		jsonb_agg(wl) AS wells
	FROM (
		SELECT jsonb_build_object(
			'well', (w.name || COALESCE(' - ' || w.well_number, '')),
			'well_total', COUNT(DISTINCT COALESCE(i.barcode, i.alt_barcode))
		) AS wl
		FROM ivs
		JOIN inventory_well AS iw ON iw.inventory_id = ANY(ivs.ids)
		JOIN inventory AS i ON i.inventory_id = iw.inventory_id
		JOIN well AS w ON w.well_id = iw.well_id
		GROUP BY w.name, w.well_number
		ORDER BY COUNT(DISTINCT COALESCE(i.barcode, i.alt_barcode))  DESC,
			w.name, w.well_number
			LIMIT 100
	) AS q
)
SELECT barcodes, barcode_total, kws.keywords, cts.container, cls.collections,
	bhs.boreholes, ocs.outcrops, sls.shotlines, ws.wells
FROM ivs, bcs, kws, cts, cls, bhs, ocs, sls, ws
LIMIT 100
