SELECT MIN(container_id) AS cid, COUNT(*) AS cid_count
FROM container
WHERE barcode = $1
