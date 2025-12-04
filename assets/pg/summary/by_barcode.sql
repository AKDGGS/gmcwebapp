WITH ivs AS (
	SELECT
		i.inventory_id,
		i.barcode,
		i.alt_barcode,
		i.keywords,
		i.collection_id,
		c.container_id,
		c.path_cache
	FROM inventory AS i
	LEFT JOIN container AS c
		ON c.container_id = i.container_id AND c.active
	WHERE i.active
	AND $1 IN (i.barcode, i.alt_barcode, c.barcode, c.alt_barcode)
),
parent_containers AS (
	SELECT array_to_json(
		ARRAY_AGG(DISTINCT ivs.path_cache ORDER BY ivs.path_cache)
	) AS parent_containers
	FROM ivs
	WHERE ivs.container_id IS NOT NULL
),
child_containers AS (
	SELECT array_to_json(
		ARRAY_AGG(DISTINCT COALESCE(c1.barcode, c1.alt_barcode)
			ORDER BY COALESCE(c1.barcode, c1.alt_barcode))
	) AS child_containers
	FROM container AS c1
	JOIN container AS c2 ON c2.container_id = c1.parent_container_id AND c2.active
	WHERE c1.active
	AND $1 IN (c2.barcode, c2.alt_barcode)
),
bcs AS (
	SELECT
		COUNT(COALESCE(barcode, alt_barcode)) AS barcode_count,
		ARRAY_AGG(DISTINCT COALESCE(barcode, alt_barcode)
			ORDER BY COALESCE(barcode, alt_barcode)) AS barcodes
	FROM ivs
	WHERE COALESCE(barcode, alt_barcode) IS NOT NULL
),
kws AS (
	SELECT ARRAY_AGG(DISTINCT keyword ORDER BY keyword) AS keywords
	FROM (
		SELECT UNNEST(keywords) AS keyword
		FROM ivs
		WHERE keywords IS NOT NULL
	) kw
),
cls AS (
	SELECT ARRAY_AGG(DISTINCT COALESCE(col.name, 'None')
		ORDER BY COALESCE(col.name, 'None')) AS collections
	FROM ivs
	LEFT JOIN collection col ON col.collection_id = ivs.collection_id
),
bhs AS (
	SELECT jsonb_agg(bh ORDER BY bh->>'prospect', bh->>'borehole') AS boreholes
	FROM (
		SELECT jsonb_build_object(
			'prospect', ps.name,
			'borehole', b.name
		) AS bh
		FROM ivs
		JOIN inventory_borehole ib ON ib.inventory_id = ivs.inventory_id
		JOIN borehole b ON b.borehole_id = ib.borehole_id
		JOIN prospect ps ON ps.prospect_id = b.prospect_id
		WHERE b.name IS NOT NULL
		GROUP BY b.name, ps.name
		ORDER BY ps.name, b.name
		LIMIT 100
	) q
),
ocs AS (
	SELECT jsonb_agg(oc ORDER BY oc->>'outcrop') AS outcrops
	FROM (
		SELECT jsonb_build_object(
			'outcrop', (o.name || COALESCE(' - ' || o.outcrop_number, ''))
		) AS oc
		FROM ivs
		LEFT JOIN inventory_outcrop io ON io.inventory_id = ivs.inventory_id
		LEFT JOIN outcrop o ON o.outcrop_id = io.outcrop_id
		WHERE o.name IS NOT NULL
		GROUP BY o.name, o.outcrop_number
		ORDER BY o.name, o.outcrop_number
		LIMIT 100
	) q
),
sls AS (
	SELECT jsonb_agg(sl ORDER BY sl->>'shotline') AS shotlines
	FROM (
		SELECT jsonb_build_object(
			'shotline', sl.name
		) AS sl
		FROM ivs
		JOIN inventory_shotpoint isp ON isp.inventory_id = ivs.inventory_id
		JOIN shotpoint sp ON sp.shotpoint_id = isp.shotpoint_id
		JOIN shotline sl ON sl.shotline_id = sp.shotline_id
		WHERE sl.name IS NOT NULL
		GROUP BY sl.name
		ORDER BY sl.name
		LIMIT 100
	) q
),
ws AS (
	SELECT jsonb_agg(wl ORDER BY wl->>'well') AS wells
	FROM (
		SELECT jsonb_build_object(
			'well', (w.name || COALESCE(' - ' || w.well_number, ''))
		) AS wl
		FROM ivs
		JOIN inventory_well iw ON iw.inventory_id = ivs.inventory_id
		JOIN well w ON w.well_id = iw.well_id
		WHERE w.name IS NOT NULL
		GROUP BY w.name, w.well_number
		ORDER BY w.name, w.well_number
		LIMIT 100
	) q
)
SELECT
	(SELECT barcodes FROM bcs) AS barcodes,
	(SELECT barcode_count FROM bcs) AS barcode_count,
	(SELECT keywords FROM kws) AS keywords,
	(SELECT parent_containers FROM parent_containers) AS parent_containers,
	(SELECT child_containers FROM child_containers) AS child_containers,
	(SELECT collections FROM cls) AS collections,
	(SELECT boreholes FROM bhs) AS boreholes,
	(SELECT outcrops FROM ocs) AS outcrops,
	(SELECT shotlines FROM sls) AS shotlines,
	(SELECT wells FROM ws) AS wells
LIMIT 100
