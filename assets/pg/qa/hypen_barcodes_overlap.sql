SELECT REGEXP_REPLACE(barcode, '-', '') AS "Barcode"
FROM inventory
WHERE active AND POSITION('-' IN barcode) <> 0

INTERSECT

SELECT REGEXP_REPLACE(barcode, '-', '') AS "Barcode"
FROM inventory
WHERE active AND POSITION('-' IN barcode) = 0
