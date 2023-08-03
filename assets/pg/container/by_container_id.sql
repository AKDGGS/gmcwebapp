SELECT container_id AS id,
name,
path_cache AS pathCache,
remark,
barcode,
alt_barcode AS altBarcode
FROM container
WHERE container_id = $1
