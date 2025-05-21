UPDATE container SET barcode = $2
WHERE barcode = $1
	OR alt_barcode = $1
