SELECT prospect_id AS "Prospect ID",
	name AS "Prospect Name"
FROM prospect
WHERE prospect_id NOT IN (
	SELECT DISTINCT prospect_id FROM borehole
)
