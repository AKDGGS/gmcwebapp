SELECT p.prospect_id, p.name AS prospect_name,
	p.alt_names AS alt_prospect_names, p.ardf_number,
	b.borehole_id, b.name, b.alt_names,
	b.is_onshore, b.completion_date,
	b.measured_depth::double precision,
	COALESCE(measured_depth_unit::text, 'ft') AS measured_depth_unit
FROM borehole AS b
LEFT OUTER JOIN prospect AS p
	ON p.prospect_id = b.prospect_id
WHERE b.borehole_id = $1
