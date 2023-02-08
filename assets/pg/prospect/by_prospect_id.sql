SELECT prospect_id AS "ID", name, alt_names AS "AltNames", ardf_number AS "ARDFNumber"
FROM prospect
WHERE prospect_id = $1
