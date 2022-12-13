SELECT prospect_id, name AS NAME, alt_names AS ALTNames, ardf_number AS ARDFNumber
FROM prospect
WHERE prospect_id = $1
