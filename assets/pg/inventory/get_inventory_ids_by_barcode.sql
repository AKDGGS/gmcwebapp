SELECT inventory_id
FROM inventory
WHERE COALESCE(barcode, alt_barcode) = $1
