SELECT prospect_id, name, alt_names, ardf_number
FROM prospect
WHERE prospect_id = $1
