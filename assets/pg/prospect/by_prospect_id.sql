SELECT prospect_id AS "ProspectID", name AS "ProspectName", alt_names AS "AltNames", ardf_number AS "ARDFNumber"
FROM prospect
WHERE prospect_id = $1
