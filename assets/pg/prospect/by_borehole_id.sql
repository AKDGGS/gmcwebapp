SELECT p.prospect_id AS "ID", p.name AS "Name", p.alt_names AS "ALTNames", p.ardf_number AS "ARDFNumber"
FROM borehole AS b
LEFT OUTER JOIN prospect AS p
	ON p.prospect_id = b.prospect_id
WHERE b.borehole_id = $1
