SELECT borehole_id AS "Borehole ID",
	p.name AS "Prospect",
	b.name AS "Borehole Name"
FROM borehole AS b
LEFT OUTER JOIN prospect AS p
	ON p.prospect_id = b.prospect_id
WHERE borehole_id NOT IN (
	SELECT DISTINCT borehole_id
	FROM inventory_borehole AS ib
	JOIN inventory AS i
		ON i.inventory_id = ib.inventory_id
	WHERE i.active
)
