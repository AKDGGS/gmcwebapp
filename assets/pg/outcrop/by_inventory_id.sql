SELECT o.outcrop_id AS id,
	o.name,
	o.outcrop_number AS number,
	o.is_onshore AS onshore,
	o.year
FROM outcrop AS o
JOIN inventory_outcrop AS io
	ON io.outcrop_id = o.outcrop_id
WHERE io.inventory_id = $1
