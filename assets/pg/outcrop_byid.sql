SELECT outcrop_id, name, outcrop_number, is_onshore
FROM outcrop
WHERE outcrop_id = $1
