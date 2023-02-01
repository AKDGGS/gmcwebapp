SELECT o.outcrop_id AS "ID", 
	o.name,
	o.outcrop_number AS "Number", 
	o.is_onshore AS "IsOnshore", 
	o.year
FROM outcrop AS o
JOIN inventory_outcrop AS io
	ON io.outcrop_id = o.outcrop_id
WHERE io.inventory_id = $1
