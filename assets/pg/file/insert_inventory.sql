INSERT INTO inventory_file(file_id, inventory_id) VALUES($1, $2)
ON CONFLICT (file_id, inventory_id) 
DO NOTHING;