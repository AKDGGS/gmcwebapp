SELECT bh.borehole_id,
	bh.name AS borehole_name,
	bh.alt_names AS alt_borehole_names,
	bh.is_onshore, bh.completion_date,

	bh.measured_depth,
	bh.measured_depth_unit::text,
	bh.elevation,
	bh.elevation_unit::text,
	ph.prospect_id,
	ph.name AS prospect_name,
	ph.alt_names AS alt_prospect_names,
	ph.ardf_number
FROM borehole AS bh
LEFT OUTER JOIN prospect AS ph
	ON ph.prospect_id = bh.prospect_id
JOIN inventory_borehole AS ib
	ON ib.borehole_id = bh.borehole_id
WHERE ib.inventory_id = $1
