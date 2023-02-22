SELECT outcrop_id as "ID", name, outcrop_number AS "Number", is_onshore AS "Onshore", year
FROM outcrop
WHERE outcrop_id = $1
