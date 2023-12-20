DELETE FROM borehole_file WHERE file_id = $1
DELETE FROM inventory_file WHERE file_id = $1
DELETE FROM outcrop_file WHERE file_id = $1
DELETE FROM prospect_file WHERE file_id = $1
