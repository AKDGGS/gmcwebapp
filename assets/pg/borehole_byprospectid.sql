SELECT borehole_id, name
FROM borehole
WHERE prospect_id = $1
ORDER BY LOWER(name)
