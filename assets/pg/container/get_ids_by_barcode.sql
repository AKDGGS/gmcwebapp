SELECT container_id
FROM container
WHERE barcode = $1
	OR barcode = ('GMC-' || $1)
	OR alt_barcode =$1
