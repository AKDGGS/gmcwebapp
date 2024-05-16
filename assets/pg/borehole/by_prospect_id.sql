SELECT
	borehole_id AS id,
	name,
	alt_names,
	completion_date
FROM borehole
WHERE prospect_id = $1
ORDER BY LOWER(name)
