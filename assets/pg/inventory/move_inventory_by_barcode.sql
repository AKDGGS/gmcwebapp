UPDATE inventory SET container_id = $1
WHERE barcode = $2
	OR barcode = ('GMC-' || $2)
	OR alt_barcode = $2
