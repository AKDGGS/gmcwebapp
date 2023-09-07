SELECT prospect_id AS id,
	name,
	alt_names AS altNames,
	ardf_number AS ARDFNumber
FROM prospect
WHERE prospect_id = $1
