SELECT outcrop_id AS id,
	name,
	outcrop_number AS number,
	is_onshore AS onshore,
	year
FROM outcrop
WHERE outcrop_id = ANY($1)
