SELECT container_id AS "Container ID",
	barcode AS "Barcode",
	alt_barcode AS "Alternate Barcode"
FROM container
WHERE active AND
	(barcode !~* '^[a-z\-\d+]+$' OR alt_barcode !~* '^[a-z\-\d+]+$')
