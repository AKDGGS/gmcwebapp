SELECT borehole_id AS id, name, alt_names AS altNames, completion_date AS completionDate
FROM borehole
WHERE prospect_id = $1
ORDER BY LOWER(name)
