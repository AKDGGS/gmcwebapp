SELECT outcrop_id as id, name, outcrop_number AS number, is_onshore AS onshore, year
FROM outcrop
WHERE outcrop_id = $1
