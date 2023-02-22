SELECT b.name, ib.borehole_id, p.name, p.prospect_id AS "ID"
FROM inventory_borehole as ib
JOIN borehole AS b ON b.borehole_id = ib.borehole_id
LEFT OUTER JOIN prospect AS p ON p.prospect_id = b.prospect_id
WHERE ib.inventory_id = $1
