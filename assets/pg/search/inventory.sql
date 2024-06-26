SELECT
	i.inventory_id AS id,
	co.name AS collection,
	i.sample_number,
	i.slide_number,
	i.box_number,
	i.set_number,
	i.core_number,
	cd.core_diameter,
	i.interval_top,
	i.interval_bottom,
	i.keywords,
	i.barcode,
	i.remark,
	i.can_publish,
	w.well,
	o.outcrop,
	b.borehole,
	sl.shotlines,
	g.geometries
FROM inventory AS i
LEFT OUTER JOIN (
	SELECT
		iw.inventory_id,
		jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
			'id', w.well_id,
			'name', w.name,
			'altnames', w.alt_names,
			'number', w.well_number,
			'api', w.api_number
		))) AS well
	FROM inventory_well AS iw
	JOIN well AS w ON w.well_id = iw.well_id
	GROUP BY iw.inventory_id
) AS w ON w.inventory_id = i.inventory_id
LEFT OUTER JOIN (
	SELECT
		io.inventory_id,
		jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
			'id', o.outcrop_id,
			'name', o.name,
			'number', o.outcrop_number,
			'year', o.year
		))) AS outcrop
	FROM inventory_outcrop AS io
	JOIN outcrop AS o ON o.outcrop_id = io.outcrop_id
	GROUP BY io.inventory_id
) AS o ON o.inventory_id = i.inventory_id
LEFT OUTER JOIN (
	SELECT
	ib.inventory_id,
	jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
		'id', b.borehole_id,
		'name', b.name,
		'prospect', jsonb_build_object(
			'id', p.prospect_id,
			'name', p.name,
			'ardf', p.ardf_number
			)
	))) AS borehole
	FROM inventory_borehole AS ib
	JOIN borehole as b ON b.borehole_id = ib.borehole_id
	LEFT OUTER JOIN prospect AS p ON p.prospect_id = b.prospect_id
	GROUP BY ib.inventory_id, b.borehole_id, p.prospect_id
) AS b on b.inventory_id = i.inventory_id
LEFT OUTER JOIN (
	SELECT 
	sp.inventory_id, 
	jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
		'id', sl.shotline_id,
		'name', sl.name,
		'year', sl.year,
		'min', sp.shotline_min,
		'max', sp.shotline_max
	))) AS shotlines
FROM(
	SELECT 
	isp.inventory_id, 
	sp.shotline_id,
	MIN(sp.shotpoint_number) AS shotline_min,
	MAX(sp.shotpoint_number) AS shotline_max
	FROM inventory_shotpoint AS isp
	JOIN shotpoint AS sp ON sp.shotpoint_id = isp.shotpoint_id
	GROUP BY isp.inventory_id, sp.shotline_id
) AS sp
JOIN shotline AS sl ON sl.shotline_id = sp.shotline_id
GROUP BY sp.inventory_id
) as sl on sl.inventory_id = i.inventory_id
LEFT OUTER JOIN (
	SELECT
		inventory_id,
		jsonb_agg(ST_ASGeoJSON(geog)::jsonb) AS geometries
	FROM inventory_geog
	WHERE ST_NPoints(geog::geometry) < 100
	GROUP BY inventory_id
) AS g ON g.inventory_id = i.inventory_id
LEFT OUTER JOIN collection AS co
	ON co.collection_id = i.collection_id
LEFT OUTER JOIN core_diameter as cd
	on cd.core_diameter_id = i.core_diameter_id
WHERE i.active
