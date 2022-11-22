SELECT outcrop_id as ID, name, outcrop_number AS OutcropNumber, is_onshore AS IsOnshore, year
FROM outcrop
WHERE outcrop_id = $1
