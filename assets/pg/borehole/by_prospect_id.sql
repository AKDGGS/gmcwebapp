SELECT borehole_id AS "ID", name, alt_names AS "AltNames", completion_date AS "CompletionDate"
FROM borehole
WHERE prospect_id = $1
ORDER BY LOWER(name)
