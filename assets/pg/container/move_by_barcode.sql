UPDATE container SET parent_container_id = $1
WHERE barcode = $2
	OR alt_barcode = $2
