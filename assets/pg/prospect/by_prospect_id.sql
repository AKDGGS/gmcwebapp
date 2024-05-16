SELECT
	prospect_id AS id,
	name,
	alt_names,
	ardf_number AS ardf
FROM prospect
WHERE prospect_id = $1
