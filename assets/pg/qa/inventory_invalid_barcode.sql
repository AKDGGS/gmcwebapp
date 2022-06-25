SELECT inventory_id AS "Inventory_ID",
	barcode AS "Barcode",
	alt_barcode AS "Alternate Barcode"
FROM inventory
WHERE active AND
	(barcode !~* '^[a-z \-\d+]+$' OR alt_barcode !~* '^[a-z \-\d+]+$')
