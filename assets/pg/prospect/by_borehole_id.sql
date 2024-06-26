SELECT p.prospect_id AS id,
	p.name,
	p.alt_names,
	p.ardf_number
FROM borehole AS b
LEFT OUTER JOIN prospect AS p
	ON p.prospect_id = b.prospect_id
WHERE b.borehole_id = $1
