SELECT MIN(container_id) AS dest_cid, COUNT(*) AS dest_cid_count,
(SELECT COUNT(*) FROM container WHERE barcode = $1 OR alt_barcode = $1) = 1 AS src_valid
FROM container
WHERE barcode = $2
	OR alt_barcode = $2
